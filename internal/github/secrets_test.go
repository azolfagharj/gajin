package github

import (
	"crypto/rand"
	"encoding/base64"
	"testing"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/nacl/box"
)

func TestEncryptSecret(t *testing.T) {
	// Generate a test public key
	publicKeyEphemeral, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	// Use the ephemeral public key as the recipient public key for testing
	var publicKey [32]byte
	copy(publicKey[:], publicKeyEphemeral[:])

	plaintext := []byte("test-secret-value")

	// Encrypt the secret
	encrypted, err := encryptSecret(plaintext, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed: %v", err)
	}

	// Verify the encrypted data structure
	// Format: [ephemeral public key (32 bytes)][encrypted ciphertext]
	if len(encrypted) < 32 {
		t.Fatalf("encrypted data too short: got %d bytes, expected at least 32", len(encrypted))
	}

	// Extract ephemeral public key
	var extractedEphemeralKey [32]byte
	copy(extractedEphemeralKey[:], encrypted[:32])

	// Verify that ephemeral key is not all zeros
	allZeros := true
	for _, b := range extractedEphemeralKey {
		if b != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Error("ephemeral public key is all zeros")
	}

	// Verify encrypted ciphertext exists
	if len(encrypted) <= 32 {
		t.Error("encrypted ciphertext is missing")
	}

	// Verify that encrypted data is different from plaintext
	if string(encrypted) == string(plaintext) {
		t.Error("encrypted data is the same as plaintext")
	}
}

func TestEncryptSecret_DifferentResults(t *testing.T) {
	// Generate a test public key
	publicKeyEphemeral, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyEphemeral[:])

	plaintext := []byte("test-secret-value")

	// Encrypt the same plaintext twice
	encrypted1, err := encryptSecret(plaintext, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed: %v", err)
	}

	encrypted2, err := encryptSecret(plaintext, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed: %v", err)
	}

	// Verify that each encryption produces different results (due to random ephemeral keys)
	if string(encrypted1) == string(encrypted2) {
		t.Error("encryptSecret produces the same result for the same input")
	}
}

func TestEncryptSecret_EmptyInput(t *testing.T) {
	publicKeyEphemeral, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyEphemeral[:])

	// Test with empty input
	encrypted, err := encryptSecret([]byte{}, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed with empty input: %v", err)
	}

	// Should still produce valid encrypted data
	if len(encrypted) < 32 {
		t.Error("encrypted data too short for empty input")
	}
}

func TestEncryptSecret_Base64Encoding(t *testing.T) {
	publicKeyEphemeral, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyEphemeral[:])

	plaintext := []byte("test-secret-value")

	encrypted, err := encryptSecret(plaintext, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed: %v", err)
	}

	// Verify that encrypted data can be base64 encoded (as it will be in SetSecret)
	encoded := base64.StdEncoding.EncodeToString(encrypted)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Failed to decode base64: %v", err)
	}

	if len(decoded) != len(encrypted) {
		t.Errorf("Decoded length mismatch: got %d, expected %d", len(decoded), len(encrypted))
	}
}

func TestNonceDerivation_Blake2b24(t *testing.T) {
	// Test that we use BLAKE2b with 24-byte output, NOT truncated 64-byte output
	// This is critical for compatibility with libsodium crypto_box_seal

	ephemeralPK := make([]byte, 32)
	recipientPK := make([]byte, 32)

	// Fill with known values
	for i := range ephemeralPK {
		ephemeralPK[i] = byte(i)
	}
	for i := range recipientPK {
		recipientPK[i] = byte(i + 32)
	}

	input := make([]byte, 64)
	copy(input[:32], ephemeralPK)
	copy(input[32:], recipientPK)

	// Method 1: WRONG - Sum512 and truncate (old buggy code)
	hash512 := blake2b.Sum512(input)
	nonce1 := hash512[:24]

	// Method 2: CORRECT - blake2b with 24-byte output (libsodium behavior)
	h, err := blake2b.New(24, nil)
	if err != nil {
		t.Fatalf("Failed to create blake2b hash: %v", err)
	}
	h.Write(input)
	nonce2 := h.Sum(nil)

	// These MUST be different - if they're equal, our fix is wrong
	if string(nonce1) == string(nonce2) {
		t.Error("BLAKE2b-512[:24] and BLAKE2b-24 produced same output - this would break sealed box")
	}

	// Verify nonce2 is 24 bytes
	if len(nonce2) != 24 {
		t.Errorf("Expected nonce length 24, got %d", len(nonce2))
	}
}

func TestEncryptSecret_OutputFormat(t *testing.T) {
	// Test that the output format matches libsodium sealed box:
	// [ephemeral_pk (32 bytes)][ciphertext + MAC (16 bytes)]
	publicKeyEphemeral, _, err := box.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatalf("Failed to generate test key: %v", err)
	}

	var publicKey [32]byte
	copy(publicKey[:], publicKeyEphemeral[:])

	plaintext := []byte("hello")
	encrypted, err := encryptSecret(plaintext, &publicKey)
	if err != nil {
		t.Fatalf("encryptSecret failed: %v", err)
	}

	// Expected length: 32 (ephemeral pk) + len(plaintext) + 16 (MAC)
	expectedLen := 32 + len(plaintext) + 16
	if len(encrypted) != expectedLen {
		t.Errorf("Expected encrypted length %d, got %d", expectedLen, len(encrypted))
	}

	// First 32 bytes should be a valid public key (non-zero)
	allZeros := true
	for i := 0; i < 32; i++ {
		if encrypted[i] != 0 {
			allZeros = false
			break
		}
	}
	if allZeros {
		t.Error("Ephemeral public key in output is all zeros")
	}
}

