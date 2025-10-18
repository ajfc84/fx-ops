package utils

import (
	"context"
)

type RootSpecMap struct {
	Projects  []string                     `yaml:"projects"`
	Pipelines map[string][]string          `yaml:"pipelines"`
	Env       map[string]map[string]string `yaml:"env"`
	Stages    map[string]StageSpec         `yaml:"stages"`
}

type SpecsData struct {
	Projects  []string
	Pipelines map[string][]string
	Env       map[string]string
	Stages    map[string]StageSpec
}

type StageSpec struct {
	Tasks []TaskSpec `yaml:"tasks"`
}

type TaskSpec struct {
	Type string `yaml:"type"`
	Exec string `yaml:"exec"`
}

var GoTaskHandlers = map[string]func(ctx context.Context) error{
	"DockerBuild": DockerBuild,
	"SpecsCommit": SpecsCommit,
}
