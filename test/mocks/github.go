package mocks

import (
	"context"
	"fmt"

	"github.com/yourusername/easy_gh_secret/internal/github"
)

// MockClient is a mock implementation of github.Client for testing.
type MockClient struct {
	PublicKeys map[string]*github.PublicKey
	Secrets    map[string]map[string]*github.SecretMetadata
	SetErrors  map[string]error
}

// NewMockClient creates a new mock GitHub client.
func NewMockClient() *MockClient {
	return &MockClient{
		PublicKeys: make(map[string]*github.PublicKey),
		Secrets:    make(map[string]map[string]*github.SecretMetadata),
		SetErrors:  make(map[string]error),
	}
}

// GetPublicKey retrieves the public key for a repository.
func (m *MockClient) GetPublicKey(ctx context.Context, owner, repo string) (*github.PublicKey, error) {
	key := fmt.Sprintf("%s/%s", owner, repo)
	if pk, ok := m.PublicKeys[key]; ok {
		return pk, nil
	}
	// Return a default public key
	return &github.PublicKey{
		KeyID: "test-key-id",
		Key:   "dGVzdC1wdWJsaWMta2V5", // base64 encoded "test-public-key"
	}, nil
}

// SetSecret sets a secret for a repository.
func (m *MockClient) SetSecret(ctx context.Context, owner, repo, name, secretValue string) error {
	key := fmt.Sprintf("%s/%s/%s", owner, repo, name)
	if err, ok := m.SetErrors[key]; ok {
		return err
	}

	// Store the secret metadata
	if m.Secrets[fmt.Sprintf("%s/%s", owner, repo)] == nil {
		m.Secrets[fmt.Sprintf("%s/%s", owner, repo)] = make(map[string]*github.SecretMetadata)
	}
	m.Secrets[fmt.Sprintf("%s/%s", owner, repo)][name] = &github.SecretMetadata{
		Name: name,
	}

	return nil
}

// GetSecret retrieves metadata about a secret.
func (m *MockClient) GetSecret(ctx context.Context, owner, repo, name string) (*github.SecretMetadata, error) {
	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if secrets, ok := m.Secrets[repoKey]; ok {
		if secret, ok := secrets[name]; ok {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("secret not found")
}

