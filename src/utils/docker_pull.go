package utils

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/rs/zerolog/log"
)

func DockerPull(cli *client.Client, ctx context.Context, imageName, imageVersion, authStr string) error {
	fullTag := fmt.Sprintf("%s:%s", imageName, imageVersion)

	_, err := cli.ImageInspect(ctx, fullTag)
	if err == nil {
		log.Info().Str("image", fullTag).Msg("Image already exists locally")
		return nil
	}

	log.Info().Str("image", fullTag).Msg("Pulling image from registry")

	resp, err := cli.ImagePull(ctx, fullTag, image.PullOptions{RegistryAuth: authStr})
	if err != nil {
		return fmt.Errorf("failed to pull image %s: %w", fullTag, err)
	}
	defer resp.Close()

	fd, isTerm := term.GetFdInfo(os.Stdout)
	if err := jsonmessage.DisplayJSONMessagesStream(resp, os.Stdout, fd, isTerm, nil); err != nil && err != io.EOF {
		log.Warn().Err(err).Msg("Failed to stream pull progress")
	}

	log.Info().Str("image", fullTag).Msg("Image pull completed successfully")
	return nil
}
