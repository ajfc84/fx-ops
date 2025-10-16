package git

import (
	"fmt"
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

func EnsureGitIdentity(envVars map[string]string) error {
	lg := log.With().Str("component", "git").Logger()

	repo, err := git.PlainOpen(envVars["CI_PROJECT_DIR"])
	if err != nil {
		lg.Error().Err(err).Msg("failed to open git repository")
		return err
	}

	cfg, err := repo.Config()
	if err != nil {
		lg.Error().Err(err).Msg("failed to read git config")
		return err
	}

	name := cfg.User.Name
	email := cfg.User.Email

	if name == "" || email == "" {
		lg.Error().Msg("missing git identity (user.name or user.email)")
		fmt.Fprintln(os.Stderr, `
Run:
  git config user.name "Your Name"
  git config user.email "you@example.com"
`)
		os.Exit(1)
	}

	lg.Info().Str("user.name", name).Str("user.email", email).Msg("git identity verified")
	return nil
}
