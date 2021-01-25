package containers

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/agrim123/gatekeeper/pkg/logger"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Image struct {
	ID string

	Reference string
}

func SearchImage(ctx context.Context, filterMap map[string]string) (*Image, error) {
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return nil, err
	}

	filterArgs := filters.NewArgs()
	for key, val := range filterMap {
		filterArgs.Add(key, val)
	}

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All:     false,
		Filters: filterArgs,
	})

	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, fmt.Errorf("No image found")
	}

	return &Image{
		ID:        images[0].ID[7:],
		Reference: filterMap["reference"],
	}, nil
}

func BuildImage(ctx context.Context, reference, dockerfile string) (*Image, error) {
	logger.Info("Building image %s", logger.Underline(reference))
	cli, err := client.NewEnvClient()
	defer cli.Close()
	if err != nil {
		return nil, err
	}

	// Create a buffer
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	// Create a filereader
	dockerFileReader, err := os.Open(dockerfile)
	if err != nil {
		return nil, err
	}

	// Read the actual Dockerfile
	readDockerFile, err := ioutil.ReadAll(dockerFileReader)
	if err != nil {
		return nil, err
	}

	// Make a TAR header for the file
	tarHeader := &tar.Header{
		Name: dockerfile,
		Size: int64(len(readDockerFile)),
	}

	// Writes the header described for the TAR file
	err = tw.WriteHeader(tarHeader)
	if err != nil {
		return nil, err
	}

	// Writes the dockerfile data to the TAR file
	_, err = tw.Write(readDockerFile)
	if err != nil {
		return nil, err
	}

	dockerFileTarReader := bytes.NewReader(buf.Bytes())

	buildOptions := types.ImageBuildOptions{
		Context:    dockerFileTarReader,
		Dockerfile: dockerfile,
		Remove:     true,
		Tags:       []string{reference},
	}

	// Build the actual image
	imageBuildResponse, err := cli.ImageBuild(
		ctx,
		dockerFileTarReader,
		buildOptions,
	)
	if err != nil {
		return nil, err
	}

	// Read the STDOUT from the build process
	defer imageBuildResponse.Body.Close()
	_, err = io.Copy(os.Stdout, imageBuildResponse.Body)
	if err != nil {
		return nil, err
	}

	logger.Success("Successfully built image %s", logger.Underline(reference))

	if image, err := SearchImage(ctx, map[string]string{"reference": reference}); err == nil {
		return image, nil
	}

	return nil, errors.New("Unable to find image")
}
