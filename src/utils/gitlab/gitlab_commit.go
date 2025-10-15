package gitlab

import (
	"encoding/json"
	"fmt"

	"github.com/rs/zerolog/log"
	gitlab "github.com/xanzy/go-gitlab"
)

func GitlabCommit(client *gitlab.Client, environment, imageVersion string, projectID int, filePath, ref, content string) error {
	if ref == "" {
		ref = "main"
	}

	yamlPath := filePath
	commitMsg := fmt.Sprintf("v%s", imageVersion)

	action := &gitlab.CommitActionOptions{
		Action:   gitlab.FileAction(gitlab.FileUpdate),
		FilePath: gitlab.String(yamlPath),
		Content:  gitlab.String(content),
	}

	commitOptions := &gitlab.CreateCommitOptions{
		Branch:        gitlab.String(ref),
		CommitMessage: gitlab.String(commitMsg),
		Actions:       []*gitlab.CommitActionOptions{action},
	}

	log.Info().Int("project_id", projectID).Str("file_path", yamlPath).Str("branch", ref).Msg("Committing update to GitLab")

	commit, _, err := client.Commits.CreateCommit(projectID, commitOptions)
	if err != nil {
		return fmt.Errorf("failed to commit changes: %w", err)
	}

	j, _ := json.MarshalIndent(commit, "", "  ")
	fmt.Println(string(j))

	log.Info().Str("commit_id", commit.ID).Str("message", commit.Message).Msg("Commit sent successfully")

	return nil
}
