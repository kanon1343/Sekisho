package oidc

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetProviderMetadata(t *testing.T) {
	mockMetadata := ProviderMetadata{
		Issuer:                "http://localhost:8080/realms/sekisho",
		AuthorizationEndpoint: "http://localhost:8080/realms/sekisho/protocol/openid-connect/auth",
		TokenEndpoint:         "http://localhost:8080/realms/sekisho/protocol/openid-connect/token",
		JwksURI:               "http://localhost:8080/realms/sekisho/protocol/openid-connect/certs",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/.well-known/openid-configuration" {
			t.Errorf("unexpected path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockMetadata)
	}))
	defer ts.Close()

	c := NewClient(ts.URL, ts.Client())
	metadata, err := c.GetProviderMetadata(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if metadata.Issuer != mockMetadata.Issuer {
		t.Errorf("expected issuer %q, got %q", mockMetadata.Issuer, metadata.Issuer)
	}
	if metadata.AuthorizationEndpoint != mockMetadata.AuthorizationEndpoint {
		t.Errorf("expected authorization endpoint %q, got %q", mockMetadata.AuthorizationEndpoint, metadata.AuthorizationEndpoint)
	}
}

func TestPKCE(t *testing.T) {
	verifier, err := GenerateCodeVerifier()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(verifier) < 43 || len(verifier) > 128 {
		t.Errorf("verifier length should be between 43 and 128, got %d", len(verifier))
	}

	challenge := GenerateCodeChallenge(verifier)
	if len(challenge) == 0 {
		t.Errorf("expected non-empty challenge")
	}

	// verifierとchallengeが異なることを確認
	if verifier == challenge {
		t.Errorf("verifier and challenge should not be same")
	}
}
