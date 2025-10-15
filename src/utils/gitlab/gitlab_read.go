package gitlab

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/rs/zerolog/log"
	gitlab "github.com/xanzy/go-gitlab"
)

func GitlabRead(client *gitlab.Client, environment, imageVersion string, projectID int, yamlPath, ref string) (string, error) {
	fileEnc := url.PathEscape(yamlPath)

	log.Info().Int("project_id", projectID).Str("file", yamlPath).Str("ref", ref).Msg("Reading file from GitLab repository")

	raw, _, err := client.RepositoryFiles.GetRawFile(
		projectID,
		fileEnc,
		&gitlab.GetRawFileOptions{Ref: gitlab.String(ref)},
	)
	if err != nil {
		return "", fmt.Errorf("failed to read file from GitLab: %w", err)
	}

	content := string(raw)

	reImage := regexp.MustCompile(`(image:\s.*:)[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+[0-9]+)?`)
	reVersion := regexp.MustCompile(`(?m)(- name:\s*(IMAGE_VERSION|DD_VERSION)\n\s*value:\s*)[0-9]+\.[0-9]+\.[0-9]+(-[a-z]+[0-9]+)?`)

	content = reImage.ReplaceAllString(content, fmt.Sprintf("${1}%s", imageVersion))
	content = reVersion.ReplaceAllString(content, fmt.Sprintf("${1}%s", imageVersion))

	log.Info().Str("image_version", imageVersion).Msg("YAML content updated in memory")

	return content, nil
}
