package containers

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

var (
	containerTimeout time.Duration = time.Second * 2
)

type Container struct {
	ID   string
	Name string

	Image string

	Cmds [][]string

	Mounts map[string]string

	containerMounts []mount.Mount

	HostConfig container.HostConfig

	FilesToCopy []string
}

func (c *Container) normalizeMounts() {
	mountBindings := make([]mount.Mount, 0)
	for src, dst := range c.Mounts {
		// Convert file to dir to mount to container
		if filesystem.IsFile(src) {
			continue
			// src = filesystem.MoveFileToDir(src)
		}

		mountBindings = append(mountBindings, mount.Mount{
			Type:   mount.TypeBind,
			Source: src,
			Target: dst,
		})
	}

	c.containerMounts = mountBindings
}

func (c *Container) Create() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	c.normalizeMounts()

	if len(c.containerMounts) > 0 {
		c.HostConfig.Mounts = c.containerMounts
	} else {
		c.HostConfig = container.HostConfig{}
	}

	ctx := context.Background()
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: c.Image,
		Cmd:   []string{"sleep", "5000"},
	}, &c.HostConfig, nil, c.Name)
	if err != nil {
		return err
	}

	c.ID = resp.ID
	if len(resp.Warnings) > 0 {
		fmt.Println("Warnings while creating the container", resp.Warnings)
	}

	return nil
}

func (c *Container) copyFiles(ctx context.Context, cli *client.Client) {
	for _, file := range c.FilesToCopy {
		dat, _ := ioutil.ReadFile(file)

		s := strings.NewReader(string(dat))

		fmt.Println(cli.CopyToContainer(ctx, c.ID, "/root", s, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}))
	}
}

func (c *Container) runStage(ctx context.Context, cli *client.Client, stage []string) error {
	fmt.Println("Running stage", stage)
	a, err := cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd:          stage,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	hijackResponse, err := cli.ContainerExecAttach(ctx, a.ID, types.ExecConfig{})
	if err != nil {
		return err
	}

	if err := cli.ContainerExecStart(ctx, a.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	b, err := ioutil.ReadAll(hijackResponse.Reader)
	if err != nil {
		return err
	}

	fmt.Printf("Output: %s\n", b)
	return nil
}

func (c *Container) runStages(ctx context.Context, cli *client.Client) error {
	for _, stage := range c.Cmds {
		c.runStage(ctx, cli, stage)
	}

	return nil
}

func (c *Container) Start(ctx context.Context) error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, c.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("Error while starting the container", err.Error())
		return err
	}

	c.copyFiles(ctx, cli)
	c.runStages(ctx, cli)

	// Wait for container to exit
	_, err = cli.ContainerWait(ctx, c.ID)

	return err
}

func (c *Container) Stop() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	// Try to stop using default timeout we are using for beast
	err = cli.ContainerStop(context.Background(), c.ID, &containerTimeout)
	if err != nil {
		return err
	}
	fmt.Println("Stopped container with ID ", c.ID)

	return nil
}

func (c *Container) Remove() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	fmt.Println("Removing container with ID ", c.ID)
	err = cli.ContainerRemove(context.Background(), c.ID, types.ContainerRemoveOptions{})

	return err
}

func (c *Container) Cleanup() {
	fmt.Println(c.Stop())
	fmt.Println(c.Remove())
}

func (c *Container) TailLogs() {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		panic(err)
	}

	stream, err := cli.ContainerLogs(context.Background(), c.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})
	defer stream.Close()
	if err != nil {
		panic(err)
	}

	logs, _ := ioutil.ReadAll(stream)
	fmt.Println(string(logs))
}
