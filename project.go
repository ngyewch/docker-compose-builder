package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/oops"
)

type Project struct {
	RepositoryMap map[string]RepositoryBuildSpec
}

type RepositoryBuildSpec struct {
	WorkingDirectory string
	Commands         []string
	ImageIdPath      string
}

func ResolveProject(dir string) (*Project, error) {
	var repositoryMap = make(map[string]RepositoryBuildSpec)
	err := doResolveProject(repositoryMap, dir)
	if err != nil {
		return nil, err
	}
	return &Project{
		RepositoryMap: repositoryMap,
	}, nil
}

func doResolveProject(repositoryMap map[string]RepositoryBuildSpec, dir string) error {
	config, err := LoadConfigurationFromFile(filepath.Join(dir, "docker-compose-builder.yml"))
	if err != nil {
		return oops.Wrapf(err, "could not load configuration from file: %w", err)
	}
	for _, include := range config.Includes {
		expandedInclude := os.ExpandEnv(include)
		if !filepath.IsAbs(expandedInclude) {
			expandedInclude, err = filepath.Rel(dir, expandedInclude)
			if err != nil {
				return err
			}
		}
		err := doResolveProject(repositoryMap, expandedInclude)
		if err != nil {
			return err
		}
	}
	for name, repository := range config.Repositories {
		repoDir := dir
		if repository.Directory != "" {
			repoDir = repository.Directory
		}
		repositoryMap[name] = RepositoryBuildSpec{
			WorkingDirectory: repoDir,
			Commands:         repository.Commands,
			ImageIdPath:      repository.ImageIdPath,
		}
	}
	return nil
}

func (p *Project) GetRepository(repository string) *RepositoryBuildSpec {
	repo, ok := p.RepositoryMap[repository]
	if !ok {
		return nil
	}
	return &repo
}
