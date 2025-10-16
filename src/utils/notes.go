package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

func ExtractNotes(changelogPath, version string) (string, error) {
	lg := log.With().
		Str("component", "extract_notes").
		Str("changelog", changelogPath).
		Str("version", version).
		Logger()

	targetVersion := strings.SplitN(version, "-", 2)[0]
	targetVersion = strings.SplitN(targetVersion, "+", 2)[0]
	versionHeader := fmt.Sprintf("### **Version %s**", targetVersion)

	file, err := os.Open(changelogPath)
	if err != nil {
		lg.Error().Err(err).Msg("failed to open changelog file")
		return "", fmt.Errorf("failed to open changelog '%s': %w", changelogPath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	inVersion := false
	inNotes := false
	var notes []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Found the target version header
		if trimmed == versionHeader {
			inVersion = true
			continue
		}

		// Skip until we reach target version
		if !inVersion {
			continue
		}

		// If we hit a new version header after our section -> stop
		if strings.HasPrefix(trimmed, "### **Version") && trimmed != versionHeader {
			break
		}

		// Start and stop block
		if trimmed == "---" {
			if !inNotes {
				inNotes = true
				continue
			} else {
				break // stop after closing this notes block
			}
		}

		if inNotes {
			notes = append(notes, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed reading changelog: %w", err)
	}

	if len(notes) == 0 {
		return "", fmt.Errorf("no notes found for version %s", version)
	}

	return strings.TrimSpace(strings.Join(notes, "\n")), nil
}
