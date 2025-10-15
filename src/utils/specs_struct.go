package utils

import (
	"fmt"
	"os"
	"reflect"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type EnvConfig struct {
	ImageVersion string `yaml:"image_version"`
	ProjectDir   string `yaml:"project_dir,omitempty"`
}

type RootSpec struct {
	Env map[string]EnvConfig `yaml:"env"`
}

func mergeStruct(environment, global EnvConfig) EnvConfig {
	eType := reflect.TypeOf(environment)
	eVal := reflect.ValueOf(environment)
	gVal := reflect.ValueOf(global)

	merged := reflect.New(eType).Elem()

	for i := 0; i < eVal.NumField(); i++ {
		envField := eVal.Field(i)
		globalField := gVal.Field(i)

		if !envField.IsZero() {
			merged.Field(i).Set(envField)
		} else {
			merged.Field(i).Set(globalField)
		}
	}

	return merged.Interface().(EnvConfig)
}

func ReadSpecsStruct(specsFile, environmentName string) (*EnvConfig, error) {
	lg := log.With().Str("component", "specs").Str("specs_file", specsFile).Logger()

	candidates := []string{
		fmt.Sprintf("%s.yaml", specsFile),
		fmt.Sprintf("%s.yml", specsFile),
	}

	var foundPath string
	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			foundPath = path
			break
		}
	}
	if foundPath == "" {
		lg.Error().Strs("candidates", candidates).Msg("spec file not found")
		return nil, fmt.Errorf("no spec file found for base '%s'", specsFile)
	}

	data, err := os.ReadFile(foundPath)
	if err != nil {
		lg.Error().Err(err).Str("path", foundPath).Msg("failed to read spec file")
		return nil, fmt.Errorf("read error '%s': %w", foundPath, err)
	}

	var spec RootSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		lg.Error().Err(err).Str("path", foundPath).Msg("yaml parse failed")
		return nil, fmt.Errorf("yaml parse error '%s': %w", foundPath, err)
	}

	global := spec.Env["global"]
	env, hasEnv := spec.Env[environmentName]

	if !hasEnv && reflect.DeepEqual(global, EnvConfig{}) {
		return nil, nil
	}

	var merged EnvConfig
	if hasEnv {
		merged = mergeStruct(env, global)
	} else {
		merged = global
	}

	logEnv(merged, lg, environmentName, foundPath)
	return &merged, nil
}

func logEnv(env EnvConfig, lg zerolog.Logger, envName, path string) {
	envLog := lg.With().Str("env", envName).Str("path", path).Logger()

	v := reflect.ValueOf(env)
	t := reflect.TypeOf(env)

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		key := field.Tag.Get("yaml")
		if key == "" {
			key = field.Name
		}
		value := v.Field(i).Interface()
		envLog.Info().Interface(key, value).Msg("spec variable loaded")
	}

	envLog.Info().Msg("Environment specs loaded successfully")
}
