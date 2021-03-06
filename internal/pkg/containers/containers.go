package containers

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/agrim123/gatekeeper/internal/constants"
	"github.com/agrim123/gatekeeper/internal/pkg/filesystem"
	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
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
	Ctx context.Context

	ID   string
	Name string

	ImageReference string

	image Image

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
			logger.Warn("Path %s not found. Not mounting", src)
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

func (c *Container) checkPrerequisite() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	if image, err := SearchImage(c.Ctx, map[string]string{"reference": constants.BaseImageName}); err != nil {
		logger.Warn("Unable to find image %s", c.ImageReference)
		image, err = BuildImage(c.Ctx, c.ImageReference, "Dockerfile")
		if err != nil {
			return err
		}
	} else {
		c.image = *image
	}

	return nil
}

func (c *Container) Create() error {
	if err := c.checkPrerequisite(); err != nil {
		return err
	}

	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	RemoveContainerIfExistsByName(c.Ctx, c.Name)

	c.normalizeMounts()

	if len(c.containerMounts) > 0 {
		c.HostConfig.Mounts = c.containerMounts
	} else {
		c.HostConfig = container.HostConfig{}
	}

	resp, err := cli.ContainerCreate(c.Ctx, &container.Config{
		Image: c.ImageReference,
		Cmd:   containerHold,
		User:  NonRootUser,
	}, &c.HostConfig, nil, &v1.Platform{
		Architecture: "amd64",
		OS:           "linux",
	}, c.Name)
	if err != nil {
		return err
	}

	c.ID = resp.ID
	if len(resp.Warnings) > 0 {
		logger.Warn("Warnings while creating the container: %v", resp.Warnings)
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
	logger.L().P(stage.Privileged).Infof("Running stage: %s, with user: %s", logger.Bold(stage.String()), logger.Underline(stage.user))
	a, err := cli.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		User:         stage.user,
		Cmd:          stage.Command,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return err
	}

	hijackResponse, err := cli.ContainerExecAttach(ctx, a.ID, types.ExecStartCheck{})
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

	if len(b) > 0 {
		logger.Info("%s %s", logger.Bold("Output:"), b)
	}

	return nil
}

func (c *Container) runStages(ctx context.Context, cli *client.Client) error {
	for _, stage := range c.Stages {
		err := c.runStage(ctx, cli, stage)
		if err != nil {
			logger.Error("Stage `%s` failed. Error: %s", logger.Bold(stage.String()), err.Error())
		} else {
			logger.Success("Stage `%s` completed", logger.Bold(stage.String()))
		}
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
		logger.Error("Error while starting the container: %s", err.Error())
		return err
	}

	c.copyFiles(ctx, cli)

	if c.Protected {
		c.Protect(ctx, cli)
	}

	c.Stages = append(c.preStages, c.Stages...)

	c.runStages(ctx, cli)

	// Wait for container to exit
	cli.ContainerWait(ctx, c.ID, container.WaitConditionNextExit)

	return err
}

func (c *Container) Stop() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	// Try to stop using default timeout we are using for beast
	err = cli.ContainerStop(c.Ctx, c.ID, &containerTimeout)
	if err != nil {
		if err != nil {
			logger.Error("Unable to stop container: %s, Error: %s", c.ID, err.Error())
		}
		return err
	}

	logger.Info("Stopped container: %s", logger.Bold(c.ID))

	return nil
}

func (c *Container) Remove() error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	logger.Info("Removing container: %s", logger.Bold(c.ID))
	err = cli.ContainerRemove(c.Ctx, c.ID, types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	})

	if err != nil {
		logger.Error("Unable to remove container: %s, Error: %s", c.ID, err.Error())
	}

	return err
}

func (c *Container) Cleanup() {
	if err := c.Stop(); err == nil {
		err = c.Remove()
	}
}

func (c *Container) TailLogs() {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		panic(err)
	}

	stream, err := cli.ContainerLogs(c.Ctx, c.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Details:    true,
	})
	if err != nil {
		logger.Error(err.Error())
		return
	}
	defer stream.Close()

	logs, _ := ioutil.ReadAll(stream)
	fmt.Println(string(logs))
}
