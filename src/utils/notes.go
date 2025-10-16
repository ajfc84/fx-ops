package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func ExtractNotes(changelogPath, version string) (string, error) {
	lg := log.With().Str("component", "extract_notes").Str("changelog", changelogPath).Str("version", version).Logger()

	targetVersion := strings.SplitN(version, "-", 2)[0]
	versionHeader := fmt.Sprintf("### **Version %s**", targetVersion)

	file, err := os.Open(changelogPath)
	if err != nil {
		return "", fmt.Errorf("failed to open changelog '%s': %w", changelogPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inVersion := false
	inBlock := false
	var notes []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == versionHeader {
			inVersion = true
			continue
		}

		if !inVersion {
			continue
		}

		if strings.HasPrefix(trimmed, "### **Version") && trimmed != versionHeader {
			break
		}

		if trimmed == "---" {
			if !inBlock {
				inBlock = true
				continue
			} else {
				break
			}
		}

		if inBlock {
			notes = append(notes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed reading changelog: %w", err)
	}

	if len(notes) == 0 {
		return "", fmt.Errorf("no notes found for version %s", version)
	}

	text := strings.TrimSpace(strings.Join(notes, "\n"))

	lg.Info().Str("version", version).Int("lines", len(notes)).Msg("notes extracted successfully")

	return text, nil
}
