package containers

import (
	"archive/tar"
	"context"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type Image struct {
	ID string
}

var _ io.Reader = (*os.File)(nil)

func BuildImageFromTar(tarPath string) (*Image, error) {
	file, err := os.Open(tarPath)
	if err != nil {
		return nil, err
	}

	content := tar.NewReader(file)

	buildOptions := types.ImageBuildOptions{
		// Tags:       []string{challengeTag},
		Remove: true,
	}

	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	imageBuildResp, err := dockerClient.ImageBuild(context.Background(), content, buildOptions)
	if err != nil {
		return nil, err
	}
	defer imageBuildResp.Body.Close()

	// TODO: return id

	return nil, err
}
