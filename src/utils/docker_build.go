package utils

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"fx-ops/myctx"
	"fx-ops/utils/env"
	"io"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
	"github.com/rs/zerolog/log"
)

func ptr(s string) *string { return &s }

func DockerBuild(
	ctx context.Context,
) error {
	cli := myctx.Get[*client.Client](ctx, myctx.DockerClient)
	envVars := myctx.Get[map[string]string](ctx, myctx.EnvVars)
	secrets := myctx.Get[map[string]string](ctx, myctx.Secrets)
	notes := myctx.Get[string](ctx, myctx.Notes)

	srcDir := fmt.Sprintf("%s/%s", envVars["SUB_PROJECT_DIR"], envVars["SRC_DIR"])
	dockerfile := fmt.Sprintf("Dockerfile.%s", envVars["CI_ENVIRONMENT_NAME"])
	imageRef := fmt.Sprintf("%s:%s", envVars["REGISTRY_IMAGE"], envVars["IMAGE_VERSION"])

	log.Info().Str("image", imageRef).Str("dockerfile", dockerfile).Str("src", srcDir).Msg("Starting Docker build")

	tarBuf, err := createTarContext(srcDir)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create Docker build context")
		return err
	}

	labels := BuildImageLabels(envVars, notes)

	buildArgs := make(map[string]*string)
	for k, v := range env.MergeEnv(envVars, secrets) {
		buildArgs[k] = ptr(v)
	}

	buildOpts := types.ImageBuildOptions{
		Tags:       []string{imageRef},
		Dockerfile: dockerfile,
		NoCache:    true,
		Remove:     true,
		Labels:     labels,
		BuildArgs:  buildArgs,
	}

	log.Info().Str("image", imageRef).Msg("Building Docker image")

	resp, err := cli.ImageBuild(ctx, tarBuf, buildOpts)
	if err != nil {
		log.Error().Err(err).Msg("Failed to start Docker build")
		return err
	}
	defer resp.Body.Close()

	fd, isTerm := term.GetFdInfo(os.Stdout)
	if err := jsonmessage.DisplayJSONMessagesStream(resp.Body, os.Stdout, fd, isTerm, nil); err != nil && err != io.EOF {
		log.Warn().Err(err).Msg("Failed to stream build progress")
	}

	log.Info().Str("image", imageRef).Msg("Docker image built successfully")

	return nil
}

func createTarContext(srcDir string) (io.Reader, error) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	defer tw.Close()

	err := filepath.Walk(srcDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			return err
		}
		hdr.Name = relPath

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = tw.Write(data)
		return err
	})
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}
