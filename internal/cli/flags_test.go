package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRepos(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single repo",
			input:    "repo1",
			expected: []string{"repo1"},
		},
		{
			name:     "multiple repos",
			input:    "repo1,repo2,repo3",
			expected: []string{"repo1", "repo2", "repo3"},
		},
		{
			name:     "repos with spaces",
			input:    "repo1, repo2 , repo3",
			expected: []string{"repo1", "repo2", "repo3"},
		},
		{
			name:     "empty repos in list",
			input:    "repo1,,repo2",
			expected: []string{"repo1", "repo2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseRepos(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

