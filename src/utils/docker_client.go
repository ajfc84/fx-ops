package utils

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func NewDockerClient() (*client.Client, context.Context, error) {
	ctx := context.Background()

	log.Info().Str("source", "env").Msg("Initializing Docker client")

	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Docker client")
		return nil, nil, fmt.Errorf("failed to create Docker client: %w", err)
	}

	ping, err := cli.Ping(ctx)
	if err != nil {
		log.Error().Err(err).Msg("Docker daemon ping failed")
		return nil, nil, fmt.Errorf("failed to ping Docker daemon: %w", err)
	}

	log.Info().Str("api_version", ping.APIVersion).Str("os_type", ping.OSType).Msg("Docker client ready")

	return cli, ctx, nil
}
