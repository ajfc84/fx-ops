package utils

import (
	"context"
	"fmt"
	"strconv"

	"fx-ops/myctx"
	"fx-ops/utils/gitlab"

	"github.com/rs/zerolog/log"
)

func SpecsCommit(ctx context.Context) error {
	log := log.With().Str("phase", "SpecsCommit").Logger()
	log.Info().Msg("Starting SpecsCommit phase")

	envVars := myctx.Get[map[string]string](ctx, myctx.EnvVars)
	secrets := myctx.Get[map[string]string](ctx, myctx.Secrets)

	client, err := gitlab.NewGitlabClient(secrets["GITLAB_TOKEN"], envVars["CI_SERVER_URL"])
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize GitLab client")
		return err
	}

	var filePath string
	switch envVars["CI_ENVIRONMENT_NAME"] {
	case "main":
		filePath = "prod"
	case "develop":
		filePath = "dev"
	default:
		filePath = envVars["CI_ENVIRONMENT_NAME"]
	}

	environment := envVars["CI_ENVIRONMENT_NAME"]
	imageVersion := envVars["LATEST_VERSION"]
	yamlPath := fmt.Sprintf("%s/%s.yaml", filePath, filePath)
	ref := "main"
	projectIDStr := envVars["CD_PROJECT_ID"]

	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		log.Error().Err(err).Str("CI_PROJECT_ID", projectIDStr).Msg("Invalid project ID format")
		return err
	}

	content, err := gitlab.GitlabRead(client, environment, imageVersion, projectID, yamlPath, ref)
	if err != nil {
		log.Error().Err(err).Msg("Failed to read GitLab file")
		return err
	}

	if err := gitlab.GitlabCommit(client, environment, imageVersion, projectID, yamlPath, ref, content); err != nil {
		log.Error().Err(err).Msg("Failed to commit updated file to GitLab")
		return err
	}

	log.Info().Msg("SpecsCommit phase completed successfully")
	return nil
}
