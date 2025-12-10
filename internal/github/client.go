package github

import (
	"context"

	"github.com/google/go-github/v57/github"
	"golang.org/x/oauth2"
)

// Client is the interface for GitHub API operations.
type Client interface {
	// GetPublicKey retrieves the public key for a repository.
	GetPublicKey(ctx context.Context, owner, repo string) (*PublicKey, error)

	// SetSecret sets a secret for a repository.
	// The secretValue should be plaintext; it will be encrypted automatically.
	SetSecret(ctx context.Context, owner, repo, name, secretValue string) error

	// GetSecret retrieves metadata about a secret (not the value, which is not accessible).
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

// GetSecret retrieves metadata about a secret.
func (c *githubClient) GetSecret(ctx context.Context, owner, repo, name string) (*SecretMetadata, error) {
	secret, _, err := c.client.Actions.GetRepoSecret(ctx, owner, repo, name)
	if err != nil {
		return nil, err
	}

	return &SecretMetadata{
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt.String(),
		UpdatedAt: secret.UpdatedAt.String(),
	}, nil
}

