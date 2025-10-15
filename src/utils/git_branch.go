package utils

import (
	git "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

func GetGitBranch() string {
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to open git repository")
	}

	headRef, err := repo.Head()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to read git HEAD")
	}

	branchName := headRef.Name().Short()
	if branchName == "" {
		log.Fatal().Msg("Failed to determine git branch name")
	}

	log.Debug().Str("git_branch", branchName).Msg("Git branch detected successfully")

	return branchName
}
