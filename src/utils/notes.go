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

	lg.Debug().Msg("reading changelog file")

	inVersionSection := false
	inNotesBlock := false
	var notes []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.TrimSpace(line) == versionHeader {
			inVersionSection = true
			continue
		}

		if inVersionSection {
			switch strings.TrimSpace(line) {
			case "---":
				if !inNotesBlock {
					inNotesBlock = true
					continue
				} else {
					inNotesBlock = false
					break
				}
			default:
				if inNotesBlock {
					notes = append(notes, line)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		lg.Error().Err(err).Msg("failed while reading changelog")
		return "", fmt.Errorf("failed to read changelog '%s': %w", changelogPath, err)
	}

	if len(notes) == 0 {
		lg.Warn().Msgf("no notes found for version %s", version)
		return "", fmt.Errorf("notes for version %s not found", version)
	}

	notesText := strings.TrimSpace(strings.Join(notes, "\n"))
	lg.Info().
		Int("lines", len(notes)).
		Msg("notes extracted successfully")

	return notesText, nil
}
