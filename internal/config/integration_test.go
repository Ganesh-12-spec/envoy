package config

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"github.com/Ganesh-12-spec/envoy/internal/crypto"
	"github.com/Ganesh-12-spec/envoy/internal/lock"
)

func TestVaultIntegration(t *testing.T) {
	// Create an isolated temporary directory.
	// Nothing from your real .envoy directory is touched.
	dir := t.TempDir()

	vaultPath := filepath.Join(dir, "vault.json")
	lockPath := filepath.Join(dir, "vault.lock")

	// ------------------------------------------------------------
	// STEP 1: Create a secret.
	// ------------------------------------------------------------

	password := []byte("test-password")
	salt := []byte("1234567890123456")

	key, err := crypto.DeriveKey(password, salt)
	if err != nil {
		t.Fatalf("deriving key: %v", err)
	}

	plaintext := []byte("integration-secret")

	ciphertext, nonce, err := crypto.Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypting secret: %v", err)
	}

	// ------------------------------------------------------------
	// STEP 2: Build a vault containing the encrypted secret.
	// ------------------------------------------------------------

	vault := crypto.Vault{
		Secrets: map[string]crypto.Secret{
			"test/api": {
				Ciphertext: base64.StdEncoding.EncodeToString(ciphertext),
				Nonce:      base64.StdEncoding.EncodeToString(nonce),
			},
		},
	}

	// ------------------------------------------------------------
	// STEP 3: Save the vault to disk.
	// ------------------------------------------------------------

	if err := SaveVault(vault, vaultPath); err != nil {
		t.Fatalf("saving vault: %v", err)
	}

	// Make sure the vault file really exists.
	if _, err := os.Stat(vaultPath); err != nil {
		t.Fatalf("vault file was not created: %v", err)
	}

	// ------------------------------------------------------------
	// STEP 4: Acquire the same type of file lock Envoy uses.
	// ------------------------------------------------------------

	fileLock, err := lock.Acquire(lockPath)
	if err != nil {
		t.Fatalf("acquiring lock: %v", err)
	}
	defer fileLock.Release()

	// ------------------------------------------------------------
	// STEP 5: Load the vault back from disk.
	// ------------------------------------------------------------

	loadedVault, err := LoadVault(vaultPath)
	if err != nil {
		t.Fatalf("loading vault: %v", err)
	}

	secret, ok := loadedVault.Secrets["test/api"]
	if !ok {
		t.Fatal("test/api was not found in vault")
	}

	// ------------------------------------------------------------
	// STEP 6: Decode the encrypted data.
	// ------------------------------------------------------------

	loadedCiphertext, err := base64.StdEncoding.DecodeString(secret.Ciphertext)
	if err != nil {
		t.Fatalf("decoding ciphertext: %v", err)
	}

	loadedNonce, err := base64.StdEncoding.DecodeString(secret.Nonce)
	if err != nil {
		t.Fatalf("decoding nonce: %v", err)
	}

	// ------------------------------------------------------------
	// STEP 7: Decrypt the secret.
	// ------------------------------------------------------------

	decrypted, err := crypto.Decrypt(
		loadedCiphertext,
		loadedNonce,
		key,
	)
	if err != nil {
		t.Fatalf("decrypting secret: %v", err)
	}

	// ------------------------------------------------------------
	// STEP 8: Verify the complete round trip.
	//
	// plaintext
	//     ↓
	// encrypt
	//     ↓
	// vault
	//     ↓
	// save to disk
	//     ↓
	// load from disk
	//     ↓
	// decrypt
	//     ↓
	// plaintext again
	// ------------------------------------------------------------

	if string(decrypted) != string(plaintext) {
		t.Fatalf(
			"decrypted secret mismatch: expected %q, got %q",
			string(plaintext),
			string(decrypted),
		)
	}
}
