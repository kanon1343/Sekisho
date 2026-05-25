package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/hkdf"
)

// Cipher provides methods for encrypting and decrypting session data
// using AES-CFB and HMAC-SHA256 for integrity.
type Cipher struct {
	encryptKey []byte
	signKey    []byte
}

// NewCipher creates a new Cipher using the provided secret.
// It derives encryption and signing keys using HKDF-SHA256.
func NewCipher(secret []byte) (*Cipher, error) {
	if len(secret) == 0 {
		return nil, errors.New("secret cannot be empty")
	}

	// Derive 32-byte encryption key
	encReader := hkdf.New(sha256.New, secret, nil, []byte("cookie-encryption"))
	encKey := make([]byte, 32)
	if _, err := io.ReadFull(encReader, encKey); err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}

	// Derive 32-byte signing key
	signReader := hkdf.New(sha256.New, secret, nil, []byte("cookie-signing"))
	signKey := make([]byte, 32)
	if _, err := io.ReadFull(signReader, signKey); err != nil {
		return nil, fmt.Errorf("failed to derive signing key: %w", err)
	}

	return &Cipher{
		encryptKey: encKey,
		signKey:    signKey,
	}, nil
}

// Encrypt encrypts the plaintext using AES-CFB and appends an HMAC-SHA256 signature.
func (c *Cipher) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	// The IV needs to be unique, but not secure. Therefore, it is common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, fmt.Errorf("failed to read random bytes for IV: %w", err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	// Append HMAC signature
	mac := hmac.New(sha256.New, c.signKey)
	mac.Write(ciphertext)
	signature := mac.Sum(nil)

	// Result: IV + Ciphertext + HMAC
	result := append(ciphertext, signature...)
	return result, nil
}

// Decrypt verifies the HMAC-SHA256 signature and decrypts the ciphertext using AES-CFB.
func (c *Cipher) Decrypt(data []byte) ([]byte, error) {
	macSize := sha256.Size
	minSize := aes.BlockSize + macSize
	if len(data) < minSize {
		return nil, errors.New("ciphertext too short")
	}

	// Split data into ciphertext and signature
	ciphertext := data[:len(data)-macSize]
	expectedMAC := data[len(data)-macSize:]

	// Verify HMAC signature
	mac := hmac.New(sha256.New, c.signKey)
	mac.Write(ciphertext)
	actualMAC := mac.Sum(nil)

	if !hmac.Equal(actualMAC, expectedMAC) {
		return nil, errors.New("invalid signature")
	}

	// Decrypt
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	iv := ciphertext[:aes.BlockSize]
	encryptedData := ciphertext[aes.BlockSize:]

	plaintext := make([]byte, len(encryptedData))
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, encryptedData)

	return plaintext, nil
}
