package utils

import (
	git "github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

func ReadGitMetadata(repoDir string) (authorName, authorEmail, revision string) {
	repo, err := git.PlainOpenWithOptions(repoDir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Warn().Err(err).Str("repo", repoDir).Msg("Failed to open Git repository")
		return "", "", ""
	}

	headRef, err := repo.Head()
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get Git HEAD")
		return "", "", ""
	}
	revision = headRef.Hash().String()[:7]

	commit, err := repo.CommitObject(headRef.Hash())
	if err != nil {
		log.Warn().Err(err).Msg("Failed to get last commit")
		return "", "", revision
	}

	return commit.Author.Name, commit.Author.Email, revision
}
