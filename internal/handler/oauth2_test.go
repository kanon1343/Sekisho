package handler

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sekisho/internal/config"
	"sekisho/internal/oidc"
	"sekisho/internal/session"
)

// MockOIDCClient is a mock for oidc.OIDCClient
type MockOIDCClient struct {
	Metadata *oidc.ProviderMetadata
	Err      error
}

func (m *MockOIDCClient) GetProviderMetadata(ctx context.Context) (*oidc.ProviderMetadata, error) {
	return m.Metadata, m.Err
}

func TestOAuth2Handler_HandleStart(t *testing.T) {
	// Setup Cipher
	secret := make([]byte, 32)
	rand.Read(secret)
	cipher, err := session.NewCipher(secret)
	if err != nil {
		t.Fatalf("failed to create cipher: %v", err)
	}

	cfg := &config.OAuth2Config{
		ClientID:    "test-client",
		RedirectURL: "https://localhost:4180/oauth2/callback",
	}

	mockOIDC := &MockOIDCClient{
		Metadata: &oidc.ProviderMetadata{
			AuthorizationEndpoint: "http://localhost:8080/auth",
		},
	}

	handler := NewOAuth2Handler(cfg, mockOIDC, cipher)

	req := httptest.NewRequest(http.MethodGet, "/oauth2/start?rd=/dashboard", nil)
	w := httptest.NewRecorder()

	handler.HandleStart(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("expected status 302, got %d", resp.StatusCode)
	}

	loc, err := resp.Location()
	if err != nil {
		t.Fatalf("failed to get location: %v", err)
	}

	if !strings.HasPrefix(loc.String(), "http://localhost:8080/auth") {
		t.Errorf("unexpected redirect location: %s", loc.String())
	}

	q := loc.Query()
	if q.Get("client_id") != "test-client" {
		t.Errorf("expected client_id test-client, got %s", q.Get("client_id"))
	}
	if q.Get("state") == "" {
		t.Error("expected state parameter")
	}
	if q.Get("nonce") == "" {
		t.Error("expected nonce parameter")
	}
	if q.Get("code_challenge") == "" {
		t.Error("expected code_challenge parameter")
	}
	if q.Get("code_challenge_method") != "S256" {
		t.Errorf("expected S256, got %s", q.Get("code_challenge_method"))
	}

	// Check Cookie
	cookies := resp.Cookies()
	var csrfCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == "_sekisho_csrf" {
			csrfCookie = c
			break
		}
	}

	if csrfCookie == nil {
		t.Fatal("expected CSRF cookie")
	}
	if !csrfCookie.HttpOnly || !csrfCookie.Secure || csrfCookie.SameSite != http.SameSiteLaxMode {
		t.Errorf("cookie attributes are incorrect: %+v", csrfCookie)
	}

	// Verify cookie can be decrypted
	decoded, err := base64.RawURLEncoding.DecodeString(csrfCookie.Value)
	if err != nil {
		t.Fatalf("failed to decode cookie: %v", err)
	}
	decrypted, err := cipher.Decrypt(decoded)
	if err != nil {
		t.Fatalf("failed to decrypt cookie: %v", err)
	}

	parts := strings.Split(string(decrypted), "|")
	if len(parts) != 4 {
		t.Fatalf("expected 4 parts in cookie data, got %d", len(parts))
	}
	if parts[0] != q.Get("state") {
		t.Errorf("state mismatch: %s != %s", parts[0], q.Get("state"))
	}
	if parts[1] != q.Get("nonce") {
		t.Errorf("nonce mismatch: %s != %s", parts[1], q.Get("nonce"))
	}
	if parts[3] != "/dashboard" {
		t.Errorf("rd mismatch: %s != /dashboard", parts[3])
	}
}
