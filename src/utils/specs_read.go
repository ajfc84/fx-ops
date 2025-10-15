package utils

import (
	"fmt"
	"fx-ops/utils/env"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

func ReadSpecs(specsFile, environmentName string) (*SpecsData, error) {
	lg := log.With().Str("component", "specs").Str("specs_file", specsFile).Logger()

	foundPath, err := GetSpecsFilename(specsFile)
	if err != nil {
		lg.Error().Err(err).Msg("specs file not found")
		return nil, fmt.Errorf("specs filename error '%s': %w", specsFile, err)
	}

	data, err := os.ReadFile(foundPath)
	if err != nil {
		lg.Error().Err(err).Msg("failed to read spec file")
		return nil, fmt.Errorf("read error '%s': %w", foundPath, err)
	}

	var spec RootSpecMap
	if err := yaml.Unmarshal(data, &spec); err != nil {
		lg.Error().Err(err).Msg("yaml parse failed")
		return nil, fmt.Errorf("yaml parse error '%s': %w", foundPath, err)
	}

	global := spec.Env["global"]
	envSpecific, hasEnv := spec.Env[environmentName]

	var merged map[string]string
	switch {
	case hasEnv:
		merged = env.MergeEnv(global, envSpecific)
	case len(global) > 0:
		merged = global
	default:
		merged = map[string]string{}
		lg.Warn().Msg("no environment-specific or global vars found")
	}

	lg.Info().Int("projects_count", len(spec.Projects)).Int("stages_count", len(spec.Stages)).Int("phases", len(spec.Phases)).Int("env_vars", len(merged)).Msg("specs loaded successfully")

	return &SpecsData{
		Projects: spec.Projects,
		Stages:   spec.Stages,
		Env:      merged,
		Phases:   spec.Phases,
	}, nil
}
