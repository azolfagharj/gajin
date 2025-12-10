package mocks

import (
	"context"
	"fmt"

	"github.com/yourusername/easy_gh_secret/internal/github"
)

// MockClient is a mock implementation of github.Client for testing.
type MockClient struct {
	PublicKeys           map[string]*github.PublicKey
	Secrets              map[string]map[string]*github.SecretMetadata
	EnvironmentSecrets   map[string]map[string]map[string]*github.SecretMetadata // repo/env/secret
	Variables            map[string]map[string]*github.VariableMetadata
	EnvironmentVariables map[string]map[string]map[string]*github.VariableMetadata // repo/env/variable
	SetErrors            map[string]error
	RepositoryIDs       map[string]int64 // owner/repo -> ID
}

// NewMockClient creates a new mock GitHub client.
func NewMockClient() *MockClient {
	return &MockClient{
		PublicKeys:           make(map[string]*github.PublicKey),
		Secrets:              make(map[string]map[string]*github.SecretMetadata),
		EnvironmentSecrets:   make(map[string]map[string]map[string]*github.SecretMetadata),
		Variables:            make(map[string]map[string]*github.VariableMetadata),
		EnvironmentVariables: make(map[string]map[string]map[string]*github.VariableMetadata),
		SetErrors:            make(map[string]error),
		RepositoryIDs:        make(map[string]int64),
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

// GetSecret retrieves metadata about a secret (legacy method).
func (m *MockClient) GetSecret(ctx context.Context, owner, repo, name string) (*github.SecretMetadata, error) {
	return m.GetRepositorySecret(ctx, owner, repo, name)
}

// SetRepositorySecret sets a repository secret.
func (m *MockClient) SetRepositorySecret(ctx context.Context, owner, repo, name, secretValue string) error {
	return m.SetSecret(ctx, owner, repo, name, secretValue)
}

// GetRepositorySecret retrieves metadata about a repository secret.
func (m *MockClient) GetRepositorySecret(ctx context.Context, owner, repo, name string) (*github.SecretMetadata, error) {
	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if secrets, ok := m.Secrets[repoKey]; ok {
		if secret, ok := secrets[name]; ok {
			return secret, nil
		}
	}
	return nil, fmt.Errorf("secret not found")
}

// GetRepositoryID retrieves the repository ID.
func (m *MockClient) GetRepositoryID(ctx context.Context, owner, repo string) (int64, error) {
	key := fmt.Sprintf("%s/%s", owner, repo)
	if id, ok := m.RepositoryIDs[key]; ok {
		return id, nil
	}
	// Return a default ID
	m.RepositoryIDs[key] = 12345
	return 12345, nil
}

// GetEnvironmentPublicKey retrieves the public key for an environment.
func (m *MockClient) GetEnvironmentPublicKey(ctx context.Context, owner, repo, environment string) (*github.PublicKey, error) {
	key := fmt.Sprintf("%s/%s/%s", owner, repo, environment)
	if pk, ok := m.PublicKeys[key]; ok {
		return pk, nil
	}
	// Return a default public key
	return &github.PublicKey{
		KeyID: "test-env-key-id",
		Key:   "dGVzdC1lbnYtcHVibGljLWtleQ==", // base64 encoded "test-env-public-key"
	}, nil
}

// SetEnvironmentSecret sets an environment secret.
func (m *MockClient) SetEnvironmentSecret(ctx context.Context, owner, repo, environment, name, secretValue string) error {
	key := fmt.Sprintf("%s/%s/%s/%s", owner, repo, environment, name)
	if err, ok := m.SetErrors[key]; ok {
		return err
	}

	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if m.EnvironmentSecrets[repoKey] == nil {
		m.EnvironmentSecrets[repoKey] = make(map[string]map[string]*github.SecretMetadata)
	}
	if m.EnvironmentSecrets[repoKey][environment] == nil {
		m.EnvironmentSecrets[repoKey][environment] = make(map[string]*github.SecretMetadata)
	}
	m.EnvironmentSecrets[repoKey][environment][name] = &github.SecretMetadata{
		Name: name,
	}

	return nil
}

// GetEnvironmentSecret retrieves metadata about an environment secret.
func (m *MockClient) GetEnvironmentSecret(ctx context.Context, owner, repo, environment, name string) (*github.SecretMetadata, error) {
	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if envSecrets, ok := m.EnvironmentSecrets[repoKey]; ok {
		if secrets, ok := envSecrets[environment]; ok {
			if secret, ok := secrets[name]; ok {
				return secret, nil
			}
		}
	}
	return nil, fmt.Errorf("environment secret not found")
}

// SetRepositoryVariable sets a repository variable.
func (m *MockClient) SetRepositoryVariable(ctx context.Context, owner, repo, name, value string) error {
	key := fmt.Sprintf("%s/%s/%s", owner, repo, name)
	if err, ok := m.SetErrors[key]; ok {
		return err
	}

	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if m.Variables[repoKey] == nil {
		m.Variables[repoKey] = make(map[string]*github.VariableMetadata)
	}
	m.Variables[repoKey][name] = &github.VariableMetadata{
		Name:  name,
		Value: value,
	}

	return nil
}

// GetRepositoryVariable retrieves a repository variable.
func (m *MockClient) GetRepositoryVariable(ctx context.Context, owner, repo, name string) (*github.VariableMetadata, error) {
	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if variables, ok := m.Variables[repoKey]; ok {
		if variable, ok := variables[name]; ok {
			return variable, nil
		}
	}
	return nil, fmt.Errorf("variable not found")
}

// SetEnvironmentVariable sets an environment variable.
func (m *MockClient) SetEnvironmentVariable(ctx context.Context, owner, repo, environment, name, value string) error {
	key := fmt.Sprintf("%s/%s/%s/%s", owner, repo, environment, name)
	if err, ok := m.SetErrors[key]; ok {
		return err
	}

	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if m.EnvironmentVariables[repoKey] == nil {
		m.EnvironmentVariables[repoKey] = make(map[string]map[string]*github.VariableMetadata)
	}
	if m.EnvironmentVariables[repoKey][environment] == nil {
		m.EnvironmentVariables[repoKey][environment] = make(map[string]*github.VariableMetadata)
	}
	m.EnvironmentVariables[repoKey][environment][name] = &github.VariableMetadata{
		Name:  name,
		Value: value,
	}

	return nil
}

// GetEnvironmentVariable retrieves an environment variable.
func (m *MockClient) GetEnvironmentVariable(ctx context.Context, owner, repo, environment, name string) (*github.VariableMetadata, error) {
	repoKey := fmt.Sprintf("%s/%s", owner, repo)
	if envVars, ok := m.EnvironmentVariables[repoKey]; ok {
		if variables, ok := envVars[environment]; ok {
			if variable, ok := variables[name]; ok {
				return variable, nil
			}
		}
	}
	return nil, fmt.Errorf("environment variable not found")
}

