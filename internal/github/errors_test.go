package github

import (
	"net/http"
	"testing"

	"github.com/google/go-github/v57/github"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvironmentNotFoundError(t *testing.T) {
	err := &EnvironmentNotFoundError{
		Owner:       "owner",
		Repo:        "repo",
		Environment: "production",
	}

	msg := err.Error()
	assert.Contains(t, msg, "environment 'production' not found")
	assert.Contains(t, msg, "owner/repo")
	assert.Contains(t, msg, "create the environment first")
}

func TestRepositoryNotFoundError(t *testing.T) {
	err := &RepositoryNotFoundError{
		Owner: "owner",
		Repo:  "repo",
	}

	msg := err.Error()
	assert.Contains(t, msg, "repository owner/repo not found")
}

func TestSecretError(t *testing.T) {
	innerErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "Bad request",
	}

	err := &SecretError{
		Type:        "repository_secret",
		Owner:        "owner",
		Repo:         "repo",
		Name:         "SECRET1",
		Err:          innerErr,
	}

	msg := err.Error()
	assert.Contains(t, msg, "failed to set repository_secret 'SECRET1'")
	assert.Contains(t, msg, "owner/repo")
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestVariableError(t *testing.T) {
	innerErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "Bad request",
	}

	err := &VariableError{
		Type:        "environment_variable",
		Owner:       "owner",
		Repo:        "repo",
		Environment: "production",
		Name:        "VAR1",
		Err:         innerErr,
	}

	msg := err.Error()
	assert.Contains(t, msg, "failed to set environment_variable 'VAR1'")
	assert.Contains(t, msg, "environment 'production'")
	assert.Contains(t, msg, "owner/repo")
	assert.Equal(t, innerErr, err.Unwrap())
}

func TestHandleGitHubError_404_Environment(t *testing.T) {
	ghErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusNotFound,
		},
		Message: "Not found",
	}

	err := handleGitHubError(ghErr, "owner", "repo", "production", "environment_secret", "SECRET1")

	envErr, ok := err.(*EnvironmentNotFoundError)
	assert.True(t, ok)
	assert.Equal(t, "owner", envErr.Owner)
	assert.Equal(t, "repo", envErr.Repo)
	assert.Equal(t, "production", envErr.Environment)
}

func TestHandleGitHubError_404_Repository(t *testing.T) {
	ghErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusNotFound,
		},
		Message: "Not found",
	}

	err := handleGitHubError(ghErr, "owner", "repo", "", "repository_secret", "SECRET1")

	repoErr, ok := err.(*RepositoryNotFoundError)
	require.True(t, ok)
	assert.Equal(t, "owner", repoErr.Owner)
	assert.Equal(t, "repo", repoErr.Repo)
}

func TestHandleGitHubError_SecretError(t *testing.T) {
	ghErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "Bad request",
	}

	err := handleGitHubError(ghErr, "owner", "repo", "", "repository_secret", "SECRET1")

	secretErr, ok := err.(*SecretError)
	require.True(t, ok)
	assert.Equal(t, "repository_secret", secretErr.Type)
	assert.Equal(t, "SECRET1", secretErr.Name)
}

func TestHandleGitHubError_VariableError(t *testing.T) {
	ghErr := &github.ErrorResponse{
		Response: &http.Response{
			StatusCode: http.StatusBadRequest,
		},
		Message: "Bad request",
	}

	err := handleGitHubError(ghErr, "owner", "repo", "production", "environment_variable", "VAR1")

	varErr, ok := err.(*VariableError)
	require.True(t, ok)
	assert.Equal(t, "environment_variable", varErr.Type)
	assert.Equal(t, "VAR1", varErr.Name)
	assert.Equal(t, "production", varErr.Environment)
}

