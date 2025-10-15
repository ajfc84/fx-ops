package utils

import (
	"fmt"
	"time"
)

func BuildImageLabels(
	envVars map[string]string,
	notes string,
) map[string]string {
	authorName, authorEmail, revision := ReadGitMetadata(envVars["CI_PROJECT_DIR"])
	authors := fmt.Sprintf("%s <%s>", authorName, authorEmail)
	created := time.Now().UTC().Format(time.RFC3339)

	return map[string]string{
		"org.opencontainers.image.title":         envVars["CI_PROJECT_NAME"],
		"org.opencontainers.image.description":   "",
		"org.opencontainers.image.url":           envVars["DOMAIN"],
		"org.opencontainers.image.source":        envVars["CI_SERVER_URL"],
		"org.opencontainers.image.version":       envVars["IMAGE_VERSION"],
		"org.opencontainers.image.revision":      revision,
		"org.opencontainers.image.created":       created,
		"org.opencontainers.image.authors":       authors,
		"org.opencontainers.image.documentation": envVars["DOMAIN"],
		"org.opencontainers.image.licenses":      "Apache-2.0",
		"org.opencontainers.image.vendor":        "Fx",
		"org.opencontainers.image.notes":         notes,
	}
}
