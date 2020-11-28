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

const (
	RootUser    = "root"
	NonRootUser = "deploy"
)

var (
	containerTimeout time.Duration = time.Second * 2
	containerHold                  = []string{"sleep", "5000"}

	// NonRootUserHomeDir defines the non root user's home directory
	NonRootUserHomeDir = "/home/" + NonRootUser
)

type Container struct {
	ID   string
	Name string

	Image string

	Stages []Stage

	preStages []Stage

	Mounts map[string]string

	containerMounts []mount.Mount

	HostConfig container.HostConfig

	FilesToCopy []string

	Protected bool
}

func (c *Container) normalizeMounts() {
	mountBindings := make([]mount.Mount, 0)
	for src, dst := range c.Mounts {
		// ignore non existent paths
		if !filesystem.DoesExists(src) {
			fmt.Println("Path " + src + " not found. Not mounting.")
			continue
		}

		// Convert file to dir to mount to container
		if filesystem.IsFile(src) {
			continue
			// src = filesystem.MoveFileToDir(src)
		}

		mountBindings = append(mountBindings, mount.Mount{
			Type:     mount.TypeBind,
			Source:   src,
			Target:   dst,
			ReadOnly: true,
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
		Cmd:   containerHold,
		User:  NonRootUser,
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

		err := cli.CopyToContainer(ctx, c.ID, NonRootUserHomeDir, s, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
		if err != nil {
			panic(err)
		}
	}
}

func (c *Container) runStage(ctx context.Context, cli *client.Client, stage Stage) error {
	fmt.Println("Running stage:", strings.Join(stage.Command, " "), "with user:", stage.user)
	a, err := cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		User:         stage.user,
		Cmd:          stage.Command,
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
	for _, stage := range c.Stages {
		c.runStage(ctx, cli, stage)
	}

	return nil
}

func (c *Container) AddPreStage(stage Stage) {
	c.preStages = append(c.preStages, stage)
}

// Protect removes shells from container to prevent attaching shell
// by user. Could find a better and more effective way to achieve this.
// Problem: User can still run basic commands (echo, ls, cat) using docker exec.
// Aim: Prevent any kind of interaction with docker container
func (c *Container) Protect(ctx context.Context, cli *client.Client) {
	stage := NewStage([]string{"rm", "-f", "/bin/sh", "/bin/bash"}, true)

	c.AddPreStage(*stage)
	// c.runStage(ctx, cli, *stage)
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

	if c.Protected {
		c.Protect(ctx, cli)
	}

	c.Stages = append(c.preStages, c.Stages...)

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
	if err := c.Stop(); err != nil {
		c.Remove()
	}
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
