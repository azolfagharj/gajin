package github

import (
	"context"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client is the interface for GitHub API operations.
type Client interface {
	// Repository Secrets
	GetPublicKey(ctx context.Context, owner, repo string) (*PublicKey, error)
	SetRepositorySecret(ctx context.Context, owner, repo, name, secretValue string) error
	GetRepositorySecret(ctx context.Context, owner, repo, name string) (*SecretMetadata, error)

	// Environment Secrets
	GetEnvironmentPublicKey(ctx context.Context, owner, repo, environment string) (*PublicKey, error)
	SetEnvironmentSecret(ctx context.Context, owner, repo, environment, name, secretValue string) error
	GetEnvironmentSecret(ctx context.Context, owner, repo, environment, name string) (*SecretMetadata, error)

	// Repository Variables
	SetRepositoryVariable(ctx context.Context, owner, repo, name, value string) error
	GetRepositoryVariable(ctx context.Context, owner, repo, name string) (*VariableMetadata, error)

	// Environment Variables
	SetEnvironmentVariable(ctx context.Context, owner, repo, environment, name, value string) error
	GetEnvironmentVariable(ctx context.Context, owner, repo, environment, name string) (*VariableMetadata, error)

	// Helper methods
	GetRepositoryID(ctx context.Context, owner, repo string) (int64, error)

	// Legacy methods (for backward compatibility during migration)
	SetSecret(ctx context.Context, owner, repo, name, secretValue string) error
	GetSecret(ctx context.Context, owner, repo, name string) (*SecretMetadata, error)
}

// PublicKey represents a GitHub repository's public key for secrets encryption.
type PublicKey struct {
	KeyID string
	Key   string
}

// SecretMetadata represents metadata about a secret.
type SecretMetadata struct {
	Name      string
	CreatedAt string
	UpdatedAt string
}

// VariableMetadata represents metadata about a variable.
type VariableMetadata struct {
	Name      string
	Value     string
	CreatedAt string
	UpdatedAt string
}

// githubClient implements the Client interface using go-github.
type githubClient struct {
	client *github.Client
}

// NewClient creates a new GitHub client.
func NewClient(token string) Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return &githubClient{
		client: github.NewClient(tc),
	}
}

// GetPublicKey retrieves the public key for a repository.
func (c *githubClient) GetPublicKey(ctx context.Context, owner, repo string) (*PublicKey, error) {
	key, _, err := c.client.Actions.GetRepoPublicKey(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	return &PublicKey{
		KeyID: key.GetKeyID(),
		Key:   key.GetKey(),
	}, nil
}

// GetRepositoryID retrieves the repository ID.
func (c *githubClient) GetRepositoryID(ctx context.Context, owner, repo string) (int64, error) {
	repository, _, err := c.client.Repositories.Get(ctx, owner, repo)
	if err != nil {
		return 0, handleGitHubError(err, owner, repo, "", "", "")
	}
	return repository.GetID(), nil
}

// SetRepositorySecret sets a secret for a repository (alias for SetSecret for backward compatibility).
func (c *githubClient) SetRepositorySecret(ctx context.Context, owner, repo, name, secretValue string) error {
	return c.SetSecret(ctx, owner, repo, name, secretValue)
}

// GetRepositorySecret retrieves metadata about a repository secret (alias for GetSecret for backward compatibility).
func (c *githubClient) GetRepositorySecret(ctx context.Context, owner, repo, name string) (*SecretMetadata, error) {
	return c.GetSecret(ctx, owner, repo, name)
}

// GetSecret retrieves metadata about a secret (legacy method).
func (c *githubClient) GetSecret(ctx context.Context, owner, repo, name string) (*SecretMetadata, error) {
	secret, _, err := c.client.Actions.GetRepoSecret(ctx, owner, repo, name)
	if err != nil {
		return nil, handleGitHubError(err, owner, repo, "", "repository_secret", name)
	}

	return &SecretMetadata{
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt.String(),
		UpdatedAt: secret.UpdatedAt.String(),
	}, nil
}

