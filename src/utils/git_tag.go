package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/rs/zerolog/log"
)

func GitTag(message string, envVars map[string]string) error {
	version := envVars["IMAGE_VERSION"]
	lg := log.With().Str("component", "git").Str("version", version).Logger()

	repo, err := git.PlainOpen(envVars["CI_PROJECT_DIR"])
	if err != nil {
		return fmt.Errorf("failed to open git repo: %w", err)
	}

	headRef, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// sig := &object.Signature{
	// 	Name:  envVars["GIT_AUTHOR_NAME"],
	// 	Email: envVars["GIT_AUTHOR_EMAIL"],
	// 	When:  time.Now(),
	// }

	_, err = repo.CreateTag(version, headRef.Hash(), &git.CreateTagOptions{
		//		Tagger:  sig,
		Message: message,
	})
	if err != nil && !strings.Contains(err.Error(), "tag already exists") {
		return fmt.Errorf("failed to create tag: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic("cannot determine home directory: " + err.Error())
	}
	keyPath := filepath.Join(home, ".ssh", "id_rsa")

	auth, err := GitSSHAuth(envVars["CI_REPOSITORY_USER"], keyPath)
	if err != nil {
		return fmt.Errorf("failed to load SSH key: %w", err)
	}

	// Allow CI/CD to override remote URL dynamically
	remoteURL := envVars["CI_REMOTE_URL"]
	if remoteURL != "" {
		lg.Info().Str("remote_url", remoteURL).Msg("using override remote URL from environment")
	}

	pushOpts := &git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/tags/%[1]s:refs/tags/%[1]s", version)),
		},
		Auth:     auth,
		Progress: os.Stdout,
		Force:    true,
	}

	if remoteURL != "" {
		pushOpts.RemoteURL = remoteURL
	}

	lg.Info().Msgf("pushing tag '%s' to remote", version)
	err = repo.Push(pushOpts)
	if err != nil {
		if err == transport.ErrEmptyRemoteRepository {
			lg.Warn().Msg("remote repository empty or unreachable, ignoring")
		} else {
			return fmt.Errorf("failed to push tag '%s': %w", version, err)
		}
	}

	lg.Info().Msgf("Tagged and pushed version: %s", version)

	return nil
}
