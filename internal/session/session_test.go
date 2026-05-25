package session

import (
	"bytes"
	"testing"
)

func TestCipher_EncryptDecrypt(t *testing.T) {
	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, err := NewCipher(secret)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	plaintext := []byte("hello world, this is a secret session data")
	ciphertext, err := cipher.Encrypt(plaintext)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if bytes.Equal(plaintext, ciphertext) {
		t.Fatal("ciphertext is equal to plaintext")
	}

	decrypted, err := cipher.Decrypt(ciphertext)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted data mismatch. got=%s, want=%s", decrypted, plaintext)
	}
}

func TestCipher_DifferentSecret(t *testing.T) {
	secret1 := []byte("secret-key-1")
	secret2 := []byte("secret-key-2")

	cipher1, _ := NewCipher(secret1)
	cipher2, _ := NewCipher(secret2)

	plaintext := []byte("some data")
	ciphertext, _ := cipher1.Encrypt(plaintext)

	// Trying to decrypt with cipher2 should fail
	_, err := cipher2.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("Expected error when decrypting with different secret, got nil")
	}
}

func TestCipher_TamperedData(t *testing.T) {
	secret := []byte("secret")
	cipher, _ := NewCipher(secret)

	plaintext := []byte("data")
	ciphertext, _ := cipher.Encrypt(plaintext)

	// Tamper with the ciphertext
	ciphertext[0] ^= 0xff

	_, err := cipher.Decrypt(ciphertext)
	if err == nil {
		t.Fatal("Expected error when decrypting tampered data, got nil")
	}
}

func TestCipher_EmptySecret(t *testing.T) {
	_, err := NewCipher([]byte{})
	if err == nil {
		t.Fatal("Expected error with empty secret")
	}
}

func TestCipher_ShortCiphertext(t *testing.T) {
	cipher, _ := NewCipher([]byte("secret"))
	_, err := cipher.Decrypt([]byte("short"))
	if err == nil {
		t.Fatal("Expected error with short ciphertext")
	}
}
