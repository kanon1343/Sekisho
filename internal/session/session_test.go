package session_test

import (
	"bytes"
	"testing"

	"sekisho/internal/session"
)

func TestNewCipher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		secret  []byte
		wantErr bool
	}{
		{"valid secret", []byte("super-secret-key-32-bytes-long!"), false},
		{"empty secret", []byte{}, true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := session.NewCipher(tt.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCipher() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCipher_EncryptDecrypt(t *testing.T) {
	t.Parallel()
	
	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, err := session.NewCipher(secret)
	if err != nil {
		t.Fatalf("NewCipher failed: %v", err)
	}

	tests := []struct {
		name      string
		plaintext []byte
	}{
		{"normal text", []byte("hello world, this is a secret session data")},
		{"empty text", []byte("")},
		{"short text", []byte("hi")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ciphertext, err := cipher.Encrypt(tt.plaintext)
			if err != nil {
				t.Fatalf("Encrypt failed: %v", err)
			}

			if len(tt.plaintext) > 0 && bytes.Equal(tt.plaintext, ciphertext) {
				t.Fatal("ciphertext is equal to plaintext")
			}

			decrypted, err := cipher.Decrypt(ciphertext)
			if err != nil {
				t.Fatalf("Decrypt failed: %v", err)
			}

			if !bytes.Equal(tt.plaintext, decrypted) {
				t.Errorf("Decrypted data mismatch. got=%s, want=%s", decrypted, tt.plaintext)
			}
		})
	}
}

func TestCipher_DecryptFailures(t *testing.T) {
	t.Parallel()

	secret1 := []byte("secret-key-1")
	cipher1, _ := session.NewCipher(secret1)
	
	secret2 := []byte("secret-key-2")
	cipher2, _ := session.NewCipher(secret2)

	plaintext := []byte("data")
	validCiphertext, _ := cipher1.Encrypt(plaintext)

	tamperedCiphertext := make([]byte, len(validCiphertext))
	copy(tamperedCiphertext, validCiphertext)
	tamperedCiphertext[0] ^= 0xff

	tests := []struct {
		name       string
		cipher     *session.Cipher
		ciphertext []byte
	}{
		{"different secret", cipher2, validCiphertext},
		{"tampered data", cipher1, tamperedCiphertext},
		{"short ciphertext", cipher1, []byte("short")},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := tt.cipher.Decrypt(tt.ciphertext)
			if err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}
