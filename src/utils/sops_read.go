package utils

import (
	"bytes"
	"fmt"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func SopsRead(filePath string) ([]byte, error) {
	lg := log.With().Str("component", "sops").Str("file", filePath).Logger()
	lg.Info().Msg("Reading SOPS file (decrypt to stdout)...")

	cmd := exec.Command("sops", "-d", filePath)

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		lg.Error().Err(err).Str("stderr", stderr.String()).Msg("Failed to read/decrypt file")
		return nil, fmt.Errorf("sops read failed for %s: %w", filePath, err)
	}

	data := out.Bytes()
	lg.Info().Int("bytes", len(data)).Msg("SOPS file read successfully")
	return data, nil
}
