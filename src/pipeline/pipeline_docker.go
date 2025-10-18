package pipeline

import (
	"context"
	"fx-ops/myctx"
	"fx-ops/utils"
	"os"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func PipelineDocker(
	ctx context.Context,
) {
	cli := myctx.Get[*client.Client](ctx, myctx.DockerClient)
	args := myctx.Get[utils.CLIArgs](ctx, myctx.Args)
	envVars := myctx.Get[map[string]string](ctx, myctx.EnvVars)

	log.Info().Str("phase", args.Stage).Str("project", args.Project).Msg("Starting pipeline in docker")

	projectName := "devops-pipeline"
	registryImage := "ajfc84/gitlab-default"
	imageVersion := envVars["LATEST_VERSION"]
	// registryAddr := "registry-1.docker.io"
	// registryUser := "ajfc84"
	// registryPasswd := utils.AskPassword("DOCKER_HUB_PASSWORD")

	authStr := ""

	if err := utils.DockerPull(cli, ctx, registryImage, imageVersion, authStr); err != nil {
		log.Fatal().Err(err).Msg("Docker pull failed")
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get working directory")
	}

	mounts := utils.BuildMounts(cli, ctx, pwd, envVars)

	if err := utils.DockerExec(
		cli,
		ctx,
		projectName,
		registryImage,
		imageVersion,
		args.Stage,
		args.Project,
		args.ExtraArgs,
		mounts,
	); err != nil {
		log.Error().Err(err).Msg("Docker exec failed")
		os.Exit(1)
	}
}
