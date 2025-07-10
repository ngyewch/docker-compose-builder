package main

import (
	"github.com/goccy/go-yaml"
	"io"
	"os"
)

type Configuration struct {
	Includes     []string                                `yaml:"includes"`
	Repositories map[string]RepositoryBuildConfiguration `yaml:"repositories"`
}

type RepositoryBuildConfiguration struct {
	Directory   string   `yaml:"directory"`
	Commands    []string `yaml:"commands"`
	ImageIdPath string   `yaml:"imageIdPath"`
}

func LoadConfiguration(r io.Reader) (*Configuration, error) {
	var configuration Configuration
	decoder := yaml.NewDecoder(r, yaml.Strict())
	err := decoder.Decode(&configuration)
	if err != nil {
		return nil, err
	}
	return &configuration, nil
}

func LoadConfigurationFromFile(path string) (*Configuration, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		_ = f.Close()
	}(f)
	return LoadConfiguration(f)
}
