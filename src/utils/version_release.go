package utils

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func ReleaseVersion() (string, error) {
	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	iter, err := repo.Tags()
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %w", err)
	}

	var stableList []*semver.Version

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := ref.Name().Short()
		if strings.Contains(tag, "-") {
			return nil
		}
		v, e := semver.NewVersion(tag)
		if e != nil {
			return nil
		}
		stableList = append(stableList, v)
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error iterating tags: %w", err)
	}

	if len(stableList) == 0 {
		return "0.0.0", nil
	}

	sort.Sort(semver.Collection(stableList))
	lastStable := stableList[len(stableList)-1]

	return lastStable.String(), nil
}
