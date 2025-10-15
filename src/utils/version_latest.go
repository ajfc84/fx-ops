package utils

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

func LatestVersion(envName string) (string, error) {
	releaseType := GetRelease(envName)

	repo, err := git.PlainOpen(".")
	if err != nil {
		return "", fmt.Errorf("failed to open repo: %w", err)
	}

	iter, err := repo.Tags()
	if err != nil {
		return "", fmt.Errorf("failed to list tags: %w", err)
	}

	var stableList []*semver.Version
	var preList []string

	err = iter.ForEach(func(ref *plumbing.Reference) error {
		tag := ref.Name().Short()
		parts := strings.Split(tag, "-")
		base := parts[0]

		v, e := semver.NewVersion(base)
		if e != nil {
			return nil
		}

		if len(parts) == 1 {
			stableList = append(stableList, v)
		} else {
			preList = append(preList, tag)
		}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("error iterating tags: %w", err)
	}

	if len(stableList) == 0 {
		if releaseType == "main" {
			return "0.0.0", nil
		}
		return fmt.Sprintf("0.0.1-%s0", releaseType), nil
	}

	sort.Sort(semver.Collection(stableList))
	lastStable := stableList[len(stableList)-1]

	// -------------------- MAIN --------------------
	if releaseType == "main" {
		return lastStable.String(), nil
	}

	// -------------------- INTERMEDIATE --------------------
	type preInfo struct {
		base *semver.Version
		n    int
		tag  string
	}
	var preByType []preInfo
	for _, tag := range preList {
		if !strings.Contains(tag, "-"+releaseType) {
			continue
		}
		parts := strings.Split(tag, "-")
		if len(parts) != 2 {
			continue
		}
		base, err := semver.NewVersion(parts[0])
		if err != nil {
			continue
		}
		nStr := strings.TrimPrefix(parts[1], releaseType)
		n, _ := strconv.Atoi(nStr)
		preByType = append(preByType, preInfo{base: base, n: n, tag: tag})
	}
	sort.Slice(preByType, func(i, j int) bool {
		if preByType[i].base.Equal(preByType[j].base) {
			return preByType[i].n < preByType[j].n
		}
		return preByType[i].base.LessThan(preByType[j].base)
	})

	if len(preByType) == 0 {
		var baseVersion string
		if len(stableList) > 0 {
			sort.Sort(semver.Collection(stableList))
			baseVersion = stableList[len(stableList)-1].String()
		} else {
			baseVersion = "0.0.0"
		}
		return fmt.Sprintf("%s-%s0", baseVersion, releaseType), nil
	}
	last := preByType[len(preByType)-1]

	return fmt.Sprintf("%s-%s%d", last.base, releaseType, last.n), nil
}
