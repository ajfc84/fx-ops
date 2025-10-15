package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func SopsEncrypt(filePath string) error {
	lg := log.With().Str("component", "sops").Str("file", filePath).Logger()
	lg.Info().Msg("Encrypting file in place with SOPS CLI...")

	cmd := exec.Command("sops", "-e", "-i", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		lg.Error().Err(err).Msg("Failed to encrypt file")
		return fmt.Errorf("sops encrypt failed for %s: %w", filePath, err)
	}

	lg.Info().Msg("File encrypted in place successfully")
	return nil
}
