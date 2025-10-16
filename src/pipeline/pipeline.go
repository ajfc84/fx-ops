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

	lg := log.With().Str("component", "pipeline").Str("stage", args.Phase).Logger()

	lg.Info().Msg("Starting pipeline")

	lg.Info().Msg("Running stages")
	var (
		isMultiProject bool
		stages         []string
		projects       []string
	)

	changelogPath := filepath.Join(envVars["CI_PROJECT_DIR"], "CHANGELOG")
	notes, err := utils.ExtractNotes(changelogPath, envVars["IMAGE_VERSION"])
	if err != nil {
		lg.Fatal().Err(err).Msg("no changelog notes found; tag message will be empty")
	}

	if args.Project != "" {
		isMultiProject = false
		stages = []string{args.Phase}
		projects = []string{args.Project}
	} else {
		isMultiProject = true
		stages = specsCfg.Stages[args.Phase]
		projects = specsCfg.Projects
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

			exists, err := ExecutePhase(ctx, stage, project)
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
			if isMultiProject {
				utils.GitTag(notes, envVars)
			}
		}

		lg.Info().Msg("stage completed successfully")
	}

	lg.Info().Msg("Pipeline finished")
}
