package main

import (
	"context"
	"fmt"

	composeV2Cli "github.com/compose-spec/compose-go/v2/cli"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	dc "github.com/fsouza/go-dockerclient"
	"github.com/samber/oops"
	"github.com/urfave/cli/v3"
)

func doBuild(ctx context.Context, cmd *cli.Command) error {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return oops.Wrapf(err, "could not create docker client")
	}
	defer func(dockerClient *client.Client) {
		_ = dockerClient.Close()
	}(dockerClient)

	project, err := ResolveProject(".")
	if err != nil {
		return oops.Wrapf(err, "could not resolve project")
	}

	builder := NewBuilder(dockerClient, project)

	dockerComposeOptions, err := composeV2Cli.NewProjectOptions([]string{"docker-compose.yml"}, composeV2Cli.WithOsEnv, composeV2Cli.WithDotEnv)
	if err != nil {
		return oops.Wrapf(err, "could not create project")
	}
	dockerComposeProject, err := dockerComposeOptions.LoadProject(ctx)
	if err != nil {
		return oops.Wrapf(err, "could not load project")
	}

	/*
		images, err := dockerClient.ImageList(ctx, image.ListOptions{
			All: true,
		})
		if err != nil {
			return err
		}
	*/

	for _, service := range dockerComposeProject.Services {
		repository, tag := dc.ParseRepositoryTag(service.Image)
		if tag == "" {
			tag = "latest"
		}
		if tag == "latest" {
			imageId, err := builder.BuildImage(ctx, repository, service.Build)
			if err != nil {
				return oops.Wrapf(err, "could not build image")
			}
			if imageId != "" {
				fmt.Printf("* imageId: %s\n", imageId)
				fmt.Println()
			}
		}
	}
	return nil
}

func findImage(images []image.Summary, repoTag string) *image.Summary {
	repository, tag := dc.ParseRepositoryTag(repoTag)
	for _, img := range images {
		for _, repoTag1 := range img.RepoTags {
			repository1, tag1 := dc.ParseRepositoryTag(repoTag1)
			if tag1 == "" {
				tag1 = "latest"
			}
			if (repository1 == repository) && (tag1 == tag) {
				return &img
			}
		}
	}
	return nil
}
