package github

import (
	"fmt"
	"net/http"

	"github.com/google/go-github/v57/github"
)

// EnvironmentNotFoundError represents an error when an environment is not found.
type EnvironmentNotFoundError struct {
	Owner       string
	Repo        string
	Environment string
}

func (e *EnvironmentNotFoundError) Error() string {
	return fmt.Sprintf("environment '%s' not found in repository %s/%s. Please create the environment first in GitHub repository settings", e.Environment, e.Owner, e.Repo)
}

// RepositoryNotFoundError represents an error when a repository is not found.
type RepositoryNotFoundError struct {
	Owner string
	Repo  string
}

func (e *RepositoryNotFoundError) Error() string {
	return fmt.Sprintf("repository %s/%s not found or access denied", e.Owner, e.Repo)
}

// SecretError represents an error related to secret operations.
type SecretError struct {
	Type        string // "repository_secret", "environment_secret"
	Owner       string
	Repo        string
	Environment string // optional, empty for repository secrets
	Name        string
	Err         error
}

func (e *SecretError) Error() string {
	if e.Environment != "" {
		return fmt.Sprintf("failed to set %s '%s' in environment '%s' for repository %s/%s: %v", e.Type, e.Name, e.Environment, e.Owner, e.Repo, e.Err)
	}
	return fmt.Sprintf("failed to set %s '%s' for repository %s/%s: %v", e.Type, e.Name, e.Owner, e.Repo, e.Err)
}

func (e *SecretError) Unwrap() error {
	return e.Err
}

// VariableError represents an error related to variable operations.
type VariableError struct {
	Type        string // "repository_variable", "environment_variable"
	Owner       string
	Repo        string
	Environment string // optional, empty for repository variables
	Name        string
	Err         error
}

func (e *VariableError) Error() string {
	if e.Environment != "" {
		return fmt.Sprintf("failed to set %s '%s' in environment '%s' for repository %s/%s: %v", e.Type, e.Name, e.Environment, e.Owner, e.Repo, e.Err)
	}
	return fmt.Sprintf("failed to set %s '%s' for repository %s/%s: %v", e.Type, e.Name, e.Owner, e.Repo, e.Err)
}

func (e *VariableError) Unwrap() error {
	return e.Err
}

// handleGitHubError converts GitHub API errors to custom error types.
func handleGitHubError(err error, owner, repo, environment, resourceType, name string) error {
	if err == nil {
		return nil
	}

	// Check if it's a GitHub API error
	ghErr, ok := err.(*github.ErrorResponse)
	if !ok {
		return err
	}

	// Handle 404 errors
	if ghErr.Response != nil && ghErr.Response.StatusCode == http.StatusNotFound {
		if environment != "" {
			return &EnvironmentNotFoundError{
				Owner:       owner,
				Repo:        repo,
				Environment: environment,
			}
		}
		return &RepositoryNotFoundError{
			Owner: owner,
			Repo:  repo,
		}
	}

	// Wrap other errors based on resource type
	if resourceType == "repository_secret" || resourceType == "environment_secret" {
		return &SecretError{
			Type:        resourceType,
			Owner:       owner,
			Repo:        repo,
			Environment: environment,
			Name:        name,
			Err:         err,
		}
	}

	if resourceType == "repository_variable" || resourceType == "environment_variable" {
		return &VariableError{
			Type:        resourceType,
			Owner:       owner,
			Repo:        repo,
			Environment: environment,
			Name:        name,
			Err:         err,
		}
	}

	return err
}

