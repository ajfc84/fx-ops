package utils

import (
	"fmt"
	"fx-ops/utils/env"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type RootSecretMap struct {
	Env map[string]map[string]string `yaml:"secrets"`
}

func ReadSecrets(secretsFile, environmentName string) (map[string]string, error) {
	lg := log.With().Str("component", "secrets").Str("secrets_file", secretsFile).Logger()

	foundPath, err := GetSpecsFilename(secretsFile)
	if err != nil {
		lg.Error().Str("secretsFile", secretsFile).Msg("secrets file not found")
		return nil, fmt.Errorf("secrets filename error '%s': %w", secretsFile, err)
	}

	data, err := SopsRead(foundPath)
	if err != nil {
		lg.Error().Err(err).Str("path", foundPath).Msg("failed to read secrets file")
		return nil, fmt.Errorf("read error '%s': %w", foundPath, err)
	}

	var spec RootSecretMap
	if err := yaml.Unmarshal(data, &spec); err != nil {
		lg.Error().Err(err).Str("path", foundPath).Msg("yaml parse failed")
		return nil, fmt.Errorf("yaml parse error '%s': %w", foundPath, err)
	}

	global := spec.Env["global"]
	environment, hasEnv := spec.Env[environmentName]

	if !hasEnv && len(global) == 0 {
		lg.Warn().Msg("no environment-specific or global vars found")
		return nil, nil
	}

	var merged map[string]string
	if hasEnv {
		merged = env.MergeEnv(global, environment)
	} else {
		merged = global
	}

	return merged, nil
}
