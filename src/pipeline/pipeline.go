package pipeline

import (
	"context"
	"fx-ops/myctx"
	"fx-ops/utils"
	"fx-ops/utils/env"
	"path/filepath"
	"strconv"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func Pipeline(
	ctx context.Context,
) {
	cli := myctx.Get[*client.Client](ctx, myctx.DockerClient)
	args := myctx.Get[utils.CLIArgs](ctx, myctx.Args)
	envVars := myctx.Get[map[string]string](ctx, myctx.EnvVars)
	secrets := myctx.Get[map[string]string](ctx, myctx.Secrets)
	specsCfg := myctx.Get[*utils.SpecsData](ctx, myctx.Config)

	utils.PrintBannerPipeline()

	lg := log.With().Str("component", "pipeline").Str("stage", args.Stage).Logger()

	lg.Info().Msg("Starting pipeline")

	isMultiProject := args.Project == ""
	stages := []string{args.Stage}
	projects := []string{args.Project}
	if isMultiProject {
		stages = specsCfg.Pipelines[args.Stage]
		if len(stages) == 0 {
			lg.Error().Str("pipeline", args.Stage).Msg("Pipeline not defined")
			utils.PrintUsage()
			return
		}
		projects = specsCfg.Projects
	}

	lg.Info().Msg("Running stages")

	changelogPath := filepath.Join(envVars["CI_PROJECT_DIR"], "CHANGELOG")
	notes, err := utils.ExtractNotes(changelogPath, envVars["IMAGE_VERSION"])
	if err != nil {
		lg.Fatal().Err(err).Msg("no changelog notes found; tag message will be empty")
	}

	envVars["IS_RELEASE"] = strconv.FormatBool(isMultiProject)

	for _, stage := range stages {
		for _, project := range projects {
			subProjectDir := filepath.Join(envVars["CI_PROJECT_DIR"], project)
			lg.Info().Msgf("%sing %s project in %s", stage, project, subProjectDir)

			lSpecsCfg, err := utils.ReadSpecs(filepath.Join(subProjectDir, "project"), envVars["CI_ENVIRONMENT_NAME"])
			if err != nil {
				lg.Fatal().Err(err).Msg("failed to load specs")
			}
			lEnvVars := env.MergeEnv(envVars, lSpecsCfg.Env)
			lEnvVars = env.ExpandVars(lEnvVars)
			lEnvVars["SUB_PROJECT_DIR"] = subProjectDir
			env.LogEnvVars(lEnvVars)

			ctx = myctx.Set(ctx, myctx.EnvVars, lEnvVars)
			ctx = myctx.Set(ctx, myctx.Config, lSpecsCfg)
			ctx = myctx.Set(ctx, myctx.Notes, notes)

			exists, err := ExecuteStage(ctx, stage, project)
			if err != nil {
				lg.Error().Err(err).Msg("Stage failed")
				return
			}

			if stage == "build" && exists {
				if isMultiProject {
					authStr, err := utils.DockerLogin(
						cli,
						ctx,
						lEnvVars["CI_REGISTRY"],
						lEnvVars["CI_REGISTRY_USER"],
						secrets["CI_REGISTRY_PASSWORD"],
					)
					if err != nil {
						lg.Fatal().Err(err).Msg("docker login failed")
					}

					if err := utils.DockerPush(
						cli,
						ctx,
						lEnvVars["CI_REGISTRY"],
						lEnvVars["REGISTRY_IMAGE"],
						lEnvVars["IMAGE_VERSION"],
						authStr,
					); err != nil {
						lg.Fatal().Err(err).Msg("docker push failed")
					}
				}
			}
		}

		if stage == "build" {
			utils.GitTag(notes, envVars)
			if isMultiProject {
				envVars["IMAGE_VERSION"], _ = utils.Version(envVars["CI_ENVIRONMENT_NAME"], false)
				envVars["LATEST_VERSION"], _ = utils.LatestVersion(envVars["CI_ENVIRONMENT_NAME"])
				envVars["RELEASE_VERSION"], _ = utils.ReleaseVersion()
			}
		}

		lg.Info().Msg("stage completed successfully")
	}

	lg.Info().Msg("Pipeline finished")
}
