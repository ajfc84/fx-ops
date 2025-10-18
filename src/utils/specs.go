package utils

import (
	"context"

	"github.com/docker/docker/client"
)

type RootSpecMap struct {
	Projects  []string                     `yaml:"projects"`
	Pipelines map[string][]string          `yaml:"pipelines"`
	Env       map[string]map[string]string `yaml:"env"`
	Stages    map[string]PhaseSpec         `yaml:"stages"`
}

type SpecsData struct {
	Projects  []string
	Pipelines map[string][]string
	Env       map[string]string
	Stages    map[string]PhaseSpec
}

type PhaseSpec struct {
	Tasks []TaskSpec `yaml:"tasks"`
}

type TaskSpec struct {
	Type string `yaml:"type"`
	Exec string `yaml:"exec"`
}

type PhaseHandler func(
	cli *client.Client,
	ctx context.Context,
	envVars map[string]string,
	srcDir string,
	notes string,
) error

var GoPhaseHandlers = map[string]func(ctx context.Context) error{
	"DockerBuild": DockerBuild,
	"SpecsCommit": SpecsCommit,
}
