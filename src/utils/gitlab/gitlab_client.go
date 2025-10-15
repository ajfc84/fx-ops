package gitlab

import (
	"fmt"

	"github.com/rs/zerolog/log"
	gitlab "github.com/xanzy/go-gitlab"
)

func NewGitlabClient(token, baseURL string) (*gitlab.Client, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(fmt.Sprintf("%s/api/v4", baseURL)))
	if err != nil {
		return nil, fmt.Errorf("failed to create GitLab client: %w", err)
	}

	log.Info().Str("base_url", baseURL).Msg("GitLab client initialized successfully")

	return client, nil
}
