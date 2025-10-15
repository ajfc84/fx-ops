package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/rs/zerolog/log"
)

func SopsDecrypt(filePath string) error {
	lg := log.With().Str("component", "sops").Str("file", filePath).Logger()
	lg.Info().Msg("Decrypting file in place with SOPS CLI...")

	cmd := exec.Command("sops", "-d", "-i", filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		lg.Error().Err(err).Msg("Failed to decrypt file")
		return fmt.Errorf("sops decrypt failed for %s: %w", filePath, err)
	}

	lg.Info().Msg("File decrypted in place successfully")
	return nil
}
