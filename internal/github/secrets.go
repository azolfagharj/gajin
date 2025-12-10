package github

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/google/go-github/v57/github"
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
		return fmt.Errorf("failed to create or update secret: %w", err)
	}

	return nil
}

// encryptSecret encrypts a secret value using NaCl box encryption.
func encryptSecret(plaintext []byte, publicKey *[32]byte) ([]byte, error) {
	// Generate ephemeral key pair
	publicKeyEphemeral, privateKeyEphemeral, err := box.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	// Encrypt the secret
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	encrypted := box.Seal(nil, plaintext, &nonce, publicKey, privateKeyEphemeral)

	// Prepend the ephemeral public key
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

// HashSecretName creates a hash of the secret name (used for some GitHub API operations).
func HashSecretName(name string) string {
	h := sha256.Sum256([]byte(name))
	return base64.StdEncoding.EncodeToString(h[:])
}

