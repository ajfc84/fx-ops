package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func Toggle(specsFile string) error {
	lg := log.With().Str("component", "sops").Str("file", specsFile).Logger()

	filePath, err := GetSpecsFilename(specsFile)
	if err != nil {
		lg.Error().Str("specsFile", specsFile).Msg("specs file not found")
		return fmt.Errorf("specs filename error '%s': %w", specsFile, err)
	}

	f, err := os.Open(filePath)
	if err != nil {
		lg.Error().Err(err).Msg("Failed to open file")
		return fmt.Errorf("failed to open %s: %w", filePath, err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	isEncrypted := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "sops:") {
			isEncrypted = true
			break
		}
	}

	if isEncrypted {
		lg.Info().Msg("Detected encrypted file — decrypting...")
		err := SopsDecrypt(filePath)
		if err != nil {
			return fmt.Errorf("decryption failed for %s: %w", filePath, err)
		}
		lg.Info().Msg("File decrypted successfully")
	} else {
		lg.Info().Msg("Detected plain file — encrypting...")
		err := SopsEncrypt(filePath)
		if err != nil {
			return fmt.Errorf("encryption failed for %s: %w", filePath, err)
		}
		lg.Info().Msg("File encrypted successfully")
	}

	return nil
}
