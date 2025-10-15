package utils

import (
	"context"
	"os"
	"path/filepath"
	"runtime"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog/log"
)

func BuildMounts(cli *client.Client, ctx context.Context, pwd string, envCfg map[string]string) []mount.Mount {
	home := os.Getenv("HOME")
	if home == "" && runtime.GOOS == "windows" {
		home = os.Getenv("USERPROFILE")
	}

	mounts := []mount.Mount{
		{
			Type:   mount.TypeBind,
			Source: pwd,
			Target: envCfg["CI_PROJECT_DIR"],
		},
		{
			Type:     mount.TypeBind,
			Source:   filepath.Join(home, ".ssh"),
			Target:   "/root/.ssh",
			ReadOnly: false,
		},
	}

	info, err := cli.Info(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot query Docker daemon info")
	}

	switch info.OSType {
	case "linux":
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeBind,
			Source: "/var/run/docker.sock",
			Target: "/var/run/docker.sock",
		})

	case "windows":
		mounts = append(mounts, mount.Mount{
			Type:   mount.TypeNamedPipe,
			Source: `\\.\pipe\docker_engine`,
			Target: `\\.\pipe\docker_engine`,
		})

	default:
		log.Fatal().Str("osType", info.OSType).Msg("Unsupported Docker daemon OS type")
	}

	return mounts
}
