package containers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

func RemoveContainerIfExistsByName(ctx context.Context, containerName string) {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return
	}
	filterMap := map[string]string{
		"name": containerName,
	}

	filterArgs := filters.NewArgs()
	for key, val := range filterMap {
		filterArgs.Add(key, val)
	}

	containers, _ := cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: filterArgs,
		All:     true,
	})

	for _, container := range containers {
		container := Container{ID: container.ID}

		container.Stop()
		container.Remove()
		break
	}
}
