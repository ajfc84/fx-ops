package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/rs/zerolog/log"
)

func DockerExec(
	cli *client.Client,
	ctx context.Context,
	projectName, registryImage, imageVersion, phase, project string,
	args []string,
	mounts []mount.Mount,
) error {
	imageRef := fmt.Sprintf("%s:%s", registryImage, imageVersion)

	env := []string{
		"CI_PROJECT_DIR=/workspaces",
		fmt.Sprintf("CI_ENVIRONMENT_NAME=%s", os.Getenv("CI_ENVIRONMENT_NAME")),
	}

	cmd := []string{phase}
	if project != "" {
		cmd = append(cmd, project)
	}
	if len(args) > 0 {
		cmd = append(cmd, args...)
	}

	cli.ContainerRemove(ctx, projectName, container.RemoveOptions{Force: true})

	log.Info().Str("image", imageRef).Str("phase", phase).Str("project", project).Msg("Starting container for exec")

	containerResp, err := cli.ContainerCreate(ctx,
		&container.Config{
			Image: imageRef,
			Cmd:   strslice.StrSlice(cmd),
			Env:   env,
			Tty:   false,
		},
		&container.HostConfig{
			Mounts:     mounts,
			AutoRemove: true,
		},
		nil, nil,
		projectName,
	)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create container")
		return err
	}

	attachResp, err := cli.ContainerAttach(ctx, containerResp.ID, container.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to attach to container")
		return err
	}
	defer attachResp.Close()

	if err := cli.ContainerStart(ctx, containerResp.ID, container.StartOptions{}); err != nil {
		log.Error().Err(err).Msg("Failed to start container")
		return err
	}

	go func() {
		if _, err := stdcopy.StdCopy(os.Stdout, os.Stderr, attachResp.Reader); err != nil {
			log.Warn().Err(err).Msg("Failed to stream container logs")
		}
	}()

	statusCh, errCh := cli.ContainerWait(ctx, containerResp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			log.Error().Err(err).Msg("Container wait error")
			return err
		}
	case status := <-statusCh:
		if status.StatusCode != 0 {
			return fmt.Errorf("container exited with status %d", status.StatusCode)
		}
	}

	log.Info().Str("image", imageRef).Msg("Container execution completed")
	return nil
}
