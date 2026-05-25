package oidc

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateCodeVerifier はRFC 7636に準拠したランダムなcode_verifierを生成します。
func GenerateCodeVerifier() (string, error) {
	// 32バイトのランダムデータを生成 (Base64URLエンコードで約43文字になる)
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random bytes for code_verifier: %w", err)
	}

	// Base64URLエンコード (パディングなし)
	verifier := base64.RawURLEncoding.EncodeToString(b)
	return verifier, nil
}

// GenerateCodeChallenge はcode_verifierからS256メソッドを用いたcode_challengeを生成します。
func GenerateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
