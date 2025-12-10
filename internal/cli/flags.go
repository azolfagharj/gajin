package cli

import (
	"strings"
)

// Flags represents all CLI flags.
type Flags struct {
	ConfigPath      string
	Token           string
	Owner           string
	Repos           string
	DryRun          bool
	ContinueOnError bool
	Verbose         bool
	ShowVersion     bool
}

// ParseRepos parses comma-separated repository names into a slice.
func ParseRepos(reposStr string) []string {
	if reposStr == "" {
		return nil
	}
	repos := strings.Split(reposStr, ",")
	result := make([]string, 0, len(repos))
	for _, repo := range repos {
		repo = strings.TrimSpace(repo)
		if repo != "" {
			result = append(result, repo)
		}
	}
	return result
}

