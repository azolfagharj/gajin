package github

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v57/github"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

// SetSecret sets a secret for a repository using GitHub's encrypted secrets API.
// The secretValue is plaintext and will be encrypted automatically.
func (c *githubClient) SetSecret(ctx context.Context, owner, repo, name, secretValue string) error {
	// Get the repository's public key
	publicKey, err := c.GetPublicKey(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to get public key: %w", err)
	}

	// Decode the public key
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey.Key)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Convert to [32]byte for nacl/box
	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], keyBytes)

	// Encrypt the secret value using NaCl box
	encrypted, err := encryptSecret([]byte(secretValue), &publicKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Create the secret
	secret := &github.EncryptedSecret{
		Name:           name,
		EncryptedValue: base64.StdEncoding.EncodeToString(encrypted),
		KeyID:          publicKey.KeyID,
	}

	_, err = c.client.Actions.CreateOrUpdateRepoSecret(ctx, owner, repo, secret)
	if err != nil {
		return handleGitHubError(err, owner, repo, "", "repository_secret", name)
	}

	return nil
}

// encryptSecret encrypts a secret value using NaCl sealed box format.
// GitHub Actions Secrets API expects LibSodium sealed box format (crypto_box_seal):
// Format: [ephemeral public key (32 bytes)][encrypted ciphertext + MAC (16 bytes)]
//
// The nonce is derived using BLAKE2b with a 24-byte output (NOT truncated from 64 bytes).
// This matches libsodium's crypto_generichash with outlen=24.
func encryptSecret(plaintext []byte, publicKey *[32]byte) ([]byte, error) {
	// Generate ephemeral key pair
	publicKeyEphemeral, privateKeyEphemeral, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Derive nonce using BLAKE2b with 24-byte output (crypto_generichash)
	// IMPORTANT: This is NOT the same as blake2b.Sum512()[:24]
	// LibSodium uses blake2b with digest_size=24, which produces different output
	nonceInput := make([]byte, 64)
	copy(nonceInput[:32], publicKeyEphemeral[:])
	copy(nonceInput[32:], publicKey[:])

	h, err := blake2b.New(24, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create blake2b hash: %w", err)
	}
	h.Write(nonceInput)
	nonceBytes := h.Sum(nil)

	var nonce [24]byte
	copy(nonce[:], nonceBytes)

	// Encrypt the secret using NaCl box
	encrypted := box.Seal(nil, plaintext, &nonce, publicKey, privateKeyEphemeral)

	// Prepend the ephemeral public key
	// Format: [ephemeral public key (32 bytes)][encrypted ciphertext]
	result := make([]byte, 32+len(encrypted))
	copy(result[:32], publicKeyEphemeral[:])
	copy(result[32:], encrypted)

	return result, nil
}

// EncryptSecretValue encrypts a plaintext secret value for a repository.
func EncryptSecretValue(ctx context.Context, client Client, owner, repo, secretValue string) ([]byte, error) {
	// Get the repository's public key
	publicKey, err := client.GetPublicKey(ctx, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get public key: %w", err)
	}

	// Decode the public key
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %w", err)
	}

	// Convert to [32]byte for nacl/box
	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], keyBytes)

	// Encrypt the secret value
	encrypted, err := encryptSecret([]byte(secretValue), &publicKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt secret: %w", err)
	}

	return encrypted, nil
}

// GetEnvironmentPublicKey retrieves the public key for an environment.
func (c *githubClient) GetEnvironmentPublicKey(ctx context.Context, owner, repo, environment string) (*PublicKey, error) {
	// Get repository ID first
	repoID, err := c.GetRepositoryID(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	key, _, err := c.client.Actions.GetEnvPublicKey(ctx, int(repoID), environment)
	if err != nil {
		return nil, handleGitHubError(err, owner, repo, environment, "environment_secret", "")
	}

	return &PublicKey{
		KeyID: key.GetKeyID(),
		Key:   key.GetKey(),
	}, nil
}

// SetEnvironmentSecret sets a secret for an environment using GitHub's encrypted secrets API.
// The secretValue is plaintext and will be encrypted automatically.
func (c *githubClient) SetEnvironmentSecret(ctx context.Context, owner, repo, environment, name, secretValue string) error {
	// Get the environment's public key
	publicKey, err := c.GetEnvironmentPublicKey(ctx, owner, repo, environment)
	if err != nil {
		return err
	}

	// Decode the public key
	keyBytes, err := base64.StdEncoding.DecodeString(publicKey.Key)
	if err != nil {
		return fmt.Errorf("failed to decode public key: %w", err)
	}

	// Convert to [32]byte for nacl/box
	var publicKeyBytes [32]byte
	copy(publicKeyBytes[:], keyBytes)

	// Encrypt the secret value using NaCl box
	encrypted, err := encryptSecret([]byte(secretValue), &publicKeyBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt secret: %w", err)
	}

	// Get repository ID
	repoID, err := c.GetRepositoryID(ctx, owner, repo)
	if err != nil {
		return err
	}

	// Create the secret
	secret := &github.EncryptedSecret{
		Name:           name,
		EncryptedValue: base64.StdEncoding.EncodeToString(encrypted),
		KeyID:          publicKey.KeyID,
	}

	_, err = c.client.Actions.CreateOrUpdateEnvSecret(ctx, int(repoID), environment, secret)
	if err != nil {
		return handleGitHubError(err, owner, repo, environment, "environment_secret", name)
	}

	return nil
}

// GetEnvironmentSecret retrieves metadata about an environment secret.
func (c *githubClient) GetEnvironmentSecret(ctx context.Context, owner, repo, environment, name string) (*SecretMetadata, error) {
	// Get repository ID
	repoID, err := c.GetRepositoryID(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	secret, _, err := c.client.Actions.GetEnvSecret(ctx, int(repoID), environment, name)
	if err != nil {
		return nil, handleGitHubError(err, owner, repo, environment, "environment_secret", name)
	}

	return &SecretMetadata{
		Name:      secret.Name,
		CreatedAt: secret.CreatedAt.String(),
		UpdatedAt: secret.UpdatedAt.String(),
	}, nil
}

// SetRepositoryVariable sets a variable for a repository.
// Variables are stored as plaintext (no encryption).
func (c *githubClient) SetRepositoryVariable(ctx context.Context, owner, repo, name, value string) error {
	variable := &github.ActionsVariable{
		Name:  name,
		Value: value,
	}

	// Try to update first, if it doesn't exist, create it
	_, err := c.client.Actions.UpdateRepoVariable(ctx, owner, repo, variable)
	if err != nil {
		// If update fails, try to create
		_, err = c.client.Actions.CreateRepoVariable(ctx, owner, repo, variable)
		if err != nil {
			return handleGitHubError(err, owner, repo, "", "repository_variable", name)
		}
	}

	return nil
}

// GetRepositoryVariable retrieves a repository variable (including its value).
func (c *githubClient) GetRepositoryVariable(ctx context.Context, owner, repo, name string) (*VariableMetadata, error) {
	variable, _, err := c.client.Actions.GetRepoVariable(ctx, owner, repo, name)
	if err != nil {
		return nil, handleGitHubError(err, owner, repo, "", "repository_variable", name)
	}

	return &VariableMetadata{
		Name:      variable.Name,
		Value:     variable.Value,
		CreatedAt: variable.CreatedAt.String(),
		UpdatedAt: variable.UpdatedAt.String(),
	}, nil
}

// SetEnvironmentVariable sets a variable for an environment.
// Variables are stored as plaintext (no encryption).
func (c *githubClient) SetEnvironmentVariable(ctx context.Context, owner, repo, environment, name, value string) error {
	// Get repository ID
	repoID, err := c.GetRepositoryID(ctx, owner, repo)
	if err != nil {
		return err
	}

	variable := &github.ActionsVariable{
		Name:  name,
		Value: value,
	}

	// Try to update first, if it doesn't exist, create it
	_, err = c.client.Actions.UpdateEnvVariable(ctx, int(repoID), environment, variable)
	if err != nil {
		// If update fails, try to create
		_, err = c.client.Actions.CreateEnvVariable(ctx, int(repoID), environment, variable)
		if err != nil {
			return handleGitHubError(err, owner, repo, environment, "environment_variable", name)
		}
	}

	return nil
}

// GetEnvironmentVariable retrieves an environment variable (including its value).
func (c *githubClient) GetEnvironmentVariable(ctx context.Context, owner, repo, environment, name string) (*VariableMetadata, error) {
	// Get repository ID
	repoID, err := c.GetRepositoryID(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	variable, _, err := c.client.Actions.GetEnvVariable(ctx, int(repoID), environment, name)
	if err != nil {
		return nil, handleGitHubError(err, owner, repo, environment, "environment_variable", name)
	}

	return &VariableMetadata{
		Name:      variable.Name,
		Value:     variable.Value,
		CreatedAt: variable.CreatedAt.String(),
		UpdatedAt: variable.UpdatedAt.String(),
	}, nil
}

