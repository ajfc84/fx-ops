package main

import (
	"context"
	"os"
	"time"

	"fx-ops/myctx"
	"fx-ops/pipeline"
	"fx-ops/utils"
	"fx-ops/utils/env"
	"fx-ops/utils/git"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
	})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	utils.PrintBannerOps()

	args := utils.ParseArgs()
	utils.CheckEnvironment()

	if err := utils.GitFetchTags(); err != nil {
		log.Warn().Err(err).Msg("Could not fetch tags from remote. Using local tags.")
	}

	project_dir, _ := os.Getwd()
	ciVars := map[string]string{
		"CI_ENVIRONMENT_NAME": utils.GetGitBranch(),
		"CI_PROJECT_DIR":      project_dir,
		"CI_PROJECT_NAME":     "fx",
	}
	version, err := utils.Version(ciVars["CI_ENVIRONMENT_NAME"], false) // TODO patch
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate version")
	}
	latest_version, err := utils.LatestVersion(ciVars["CI_ENVIRONMENT_NAME"])
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate latest version")
	}
	release_version, err := utils.ReleaseVersion()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate release version")
	}
	opsVars := map[string]string{
		"IMAGE_VERSION":   version,
		"LATEST_VERSION":  latest_version,
		"RELEASE_VERSION": release_version,
		"SECRETS_FILE":    "secrets",
	}
	envVars := env.MergeEnv(opsVars, ciVars)
	specsCfg, err := utils.ReadSpecs("main", envVars["CI_ENVIRONMENT_NAME"])
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load specs")
	}
	envVars = env.MergeEnv(envVars, specsCfg.Env)
	envVars = env.ExpandVars(envVars)
	env.LogEnvVars(envVars)

	if err := git.EnsureGitIdentity(envVars); err != nil {
		os.Exit(1)
	}

	if args.Install {
		log.Warn().Msg("Installing dependencies (TODO: implement buildSh/install.sh equivalent)")
	}

	switch args.Stage {
	case "version":
		version := envVars["IMAGE_VERSION"]
		log.Info().Str("IMAGE_VERSION", version).Msg("Version completed successfully")
	case "sops":
		log.Info().Msg("Running SOPS")
		if err := utils.Toggle(envVars["SECRETS_FILE"]); err != nil {
			log.Error().Err(err).Msg("SOPS toggle failed")
		} else {
			log.Info().Msg("SOPS toggle completed successfully")
		}
	default:
		secrets, err := utils.ReadSecrets(envVars["SECRETS_FILE"], envVars["CI_ENVIRONMENT_NAME"])
		if err != nil {
			log.Fatal().Err(err).Msg("failed to load secrets")
		}
		secrets["CI_REGISTRY_PASSWORD"] = secrets["REGISTRY_TOKEN"]
		env.LogSecrets(secrets)

		cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create Docker client")
		}
		defer cli.Close()

		ctx := context.Background()
		ctx = myctx.Set(ctx, myctx.DockerClient, cli)
		ctx = myctx.Set(ctx, myctx.Args, args)
		ctx = myctx.Set(ctx, myctx.EnvVars, envVars)
		ctx = myctx.Set(ctx, myctx.Secrets, secrets)
		ctx = myctx.Set(ctx, myctx.Config, specsCfg)

		if args.Docker {
			pipeline.PipelineDocker(ctx)
		} else {
			pipeline.Pipeline(ctx)
		}
	}
}
