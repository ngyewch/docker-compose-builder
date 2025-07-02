package main

import (
	"github.com/goccy/go-yaml"
	"io"
)

type Project struct {
	Repositories map[string]Repository `json:"repositories"`
}

type Repository struct {
	Directory   string   `yaml:"directory"`
	Commands    []string `yaml:"commands"`
	ImageIdPath string   `yaml:"imageIdPath"`
}

func LoadProject(r io.Reader) (*Project, error) {
	var project Project
	decoder := yaml.NewDecoder(r, yaml.Strict())
	err := decoder.Decode(&project)
	if err != nil {
		return nil, err
	}
	return &project, nil
}

func (p *Project) FindRepository(repository string) *Repository {
	repo, ok := p.Repositories[repository]
	if !ok {
		return nil
	}
	return &repo
}
