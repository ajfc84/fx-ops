package utils

import (
	"fmt"
	"os"
	"path/filepath"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/rs/zerolog/log"
)

func GitFetchTags() error {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("remote 'origin' not found: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory: " + err.Error())
	}
	keyPath := filepath.Join(home, ".ssh", "id_rsa")

	auth, err := GitSSHAuth("git", keyPath)
	if err != nil {
		return fmt.Errorf("failed to prepare SSH auth: %w", err)
	}

	log.Info().Msg("Fetching tags from origin via SSH")

	err = remote.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{
			"+refs/tags/*:refs/tags/*",
		},
		Tags:     git.AllTags,
		Force:    true,
		Prune:    true,
		Auth:     auth,
		Progress: os.Stdout,
	})

	switch err {
	case nil:
		log.Info().Msg("Tags fetched and synchronized successfully")
	case git.NoErrAlreadyUpToDate:
		log.Info().Msg("Tags already up to date")
	default:
		return fmt.Errorf("fetch tags failed: %w", err)
	}

	return nil
}
