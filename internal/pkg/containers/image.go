package containers

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Image struct {
	ID string
}

func CheckIfImageExists(ctx context.Context, reference string) error {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return err
	}

	filterMap := make(map[string]string)
	filterMap["reference"] = reference

	filterArgs := filters.NewArgs()
	for key, val := range filterMap {
		filterArgs.Add(key, val)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All:     false,
		Filters: filterArgs,
	})

	if err != nil {
		return err
	}

	if len(images) == 0 {
		return fmt.Errorf("No image %s found", reference)
	}

	return nil
}

func BuildImage(dockerfile string) error {
	return nil
}
