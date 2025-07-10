package main

import (
	"context"
	"fmt"
	"github.com/compose-spec/compose-go/v2/types"
	"github.com/docker/docker/client"
	"github.com/google/shlex"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

type Builder struct {
	dockerClient *client.Client
	project      *Project
}

func NewBuilder(dockerClient *client.Client, project *Project) *Builder {
	return &Builder{
		dockerClient: dockerClient,
		project:      project,
	}
}

func (builder *Builder) BuildImage(ctx context.Context, repository string, build *types.BuildConfig) (string, error) {
	var imageId string
	repo := builder.project.GetRepository(repository)
	if repo == nil {
		return imageId, nil
	}
	workingDirectory := os.ExpandEnv(repo.WorkingDirectory)
	for _, command := range repo.Commands {
		argv, err := shlex.Split(command)
		if err != nil {
			return imageId, err
		}
		fmt.Printf("========== Building image for %s ==========\n", repository)
		cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
		cmd.Dir = workingDirectory
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		fmt.Println()
		if err != nil {
			return imageId, err
		}
	}
	if repo.ImageIdPath != "" {
		f, err := os.Open(filepath.Join(workingDirectory, repo.ImageIdPath))
		if err != nil {
			return imageId, err
		}
		imageIdBytes, err := io.ReadAll(f)
		if err != nil {
			return imageId, err
		}
		imageId = string(imageIdBytes)
	}
	if imageId != "" {
		inspectResponse, err := builder.dockerClient.ImageInspect(ctx, imageId)
		if err != nil {
			return imageId, err
		}
		matched := false
		for _, repoTag := range inspectResponse.RepoTags {
			if repoTag == repository+":latest" {
				matched = true
				break
			}
		}
		if !matched {
			return imageId, fmt.Errorf("image %s is not tagged as '%s:latest'", imageId, repository)
		}
	}

	return imageId, nil
}
