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

func DockerPush(cli *client.Client, ctx context.Context, registry, imageName, imageVersion, authStr string) error {
	sourceTag := fmt.Sprintf("%s:%s", imageName, imageVersion)
	fullTag := sourceTag

	if registry != "registry-1.docker.io" {
		fullTag = fmt.Sprintf("%s/%s:%s", registry, imageName, imageVersion)
		if err := cli.ImageTag(ctx, sourceTag, fullTag); err != nil {
			return fmt.Errorf("failed to tag image: %w", err)
		}
	}

	log.Info().Str("image", fullTag).Msg("Pushing image to registry")

	resp, err := cli.ImagePush(ctx, fullTag, image.PushOptions{RegistryAuth: authStr})
	if err != nil {
		return fmt.Errorf("failed to push image %s: %w", fullTag, err)
	}
	defer resp.Close()

	fd, isTerm := term.GetFdInfo(os.Stdout)
	if err := jsonmessage.DisplayJSONMessagesStream(resp, os.Stdout, fd, isTerm, nil); err != nil && err != io.EOF {
		log.Warn().Err(err).Msg("Failed to stream push progress")
	}

	log.Info().Str("image", fullTag).Msg("Image push completed successfully")
	return nil
}
