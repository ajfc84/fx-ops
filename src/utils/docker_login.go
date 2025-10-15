package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func DockerLogin(cli *client.Client, ctx context.Context, registryAddr, user, password string) (string, error) {
	authConfig := registry.AuthConfig{
		ServerAddress: registryAddr,
		Username:      user,
		Password:      password,
	}

	resp, err := cli.RegistryLogin(ctx, authConfig)
	if err != nil {
		log.Error().Str("registry", registryAddr).Err(err).Msg("Docker registry login failed")
		return "", fmt.Errorf("docker login failed: %w", err)
	}

	if resp.IdentityToken != "" {
		authConfig.IdentityToken = resp.IdentityToken
		authConfig.Password = ""
	}

	authBytes, err := json.Marshal(authConfig)
	if err != nil {
		return "", fmt.Errorf("failed to marshal auth config: %w", err)
	}
	authStr := base64.URLEncoding.EncodeToString(authBytes)

	log.Info().Str("registry", registryAddr).Str("status", resp.Status).Bool("token_used", resp.IdentityToken != "").Msg("Docker registry login successful")

	return authStr, nil
}
