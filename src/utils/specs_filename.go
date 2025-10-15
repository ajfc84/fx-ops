package utils

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

func GetSpecsFilename(specsFile string) (string, error) {
	lg := log.With().
		Str("component", "specs").
		Str("specs_base", specsFile).
		Logger()

	candidates := []string{
		fmt.Sprintf("%s.yaml", specsFile),
		fmt.Sprintf("%s.yml", specsFile),
	}

	for _, path := range candidates {
		if _, err := os.Stat(path); err == nil {
			lg.Info().Str("path", path).Msg("Specs file resolved")
			return path, nil
		}
	}

	lg.Error().Strs("candidates", candidates).Msg("No specs file found")
	return "", fmt.Errorf("no specs file found for base '%s'", specsFile)
}
