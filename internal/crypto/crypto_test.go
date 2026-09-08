package crypto

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("hello secret")

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	decrypted, err := Decrypt(ciphertext, nonce, key)
	if err != nil {
		t.Fatalf("Decrypt() failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Fatalf("decrypted text = %q, want %q", decrypted, plaintext)
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key := make([]byte, 32)
	wrongKey := make([]byte, 32)

	key[0] = 1
	wrongKey[0] = 2

	plaintext := []byte("hello secret")

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	_, err = Decrypt(ciphertext, nonce, wrongKey)
	if err == nil {
		t.Fatal("Decrypt() succeeded with wrong key")
	}
}

func TestDecryptCorruptedCiphertext(t *testing.T) {
	key := make([]byte, 32)
	plaintext := []byte("hello secret")

	ciphertext, nonce, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("Encrypt() failed: %v", err)
	}

	ciphertext[0] ^= 1

	_, err = Decrypt(ciphertext, nonce, key)
	if err == nil {
		t.Fatal("Decrypt() succeeded with corrupted ciphertext")
	}
}
