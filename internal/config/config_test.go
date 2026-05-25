package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		env         map[string]string
		yamlContent string
		wantErr     bool
		check       func(t *testing.T, cfg *Config)
	}{
		{
			name: "Default and Env",
			env: map[string]string{
				"OAUTH2_CLIENT_SECRET":  "dummy-oauth2-secret",
				"SESSION_COOKIE_SECRET": "dummy-session-secret",
				"SERVER_PORT":           "8080",
			},
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != "8080" {
					t.Errorf("expected port 8080, got %s", cfg.Server.Port)
				}
				if cfg.OAuth2.ClientSecret != "dummy-oauth2-secret" {
					t.Errorf("expected secret dummy-oauth2-secret, got %s", cfg.OAuth2.ClientSecret)
				}
			},
		},
		{
			name: "Yaml and Env Override",
			env: map[string]string{
				"OAUTH2_CLIENT_SECRET":  "dummy-oauth2-secret",
				"SESSION_COOKIE_SECRET": "dummy-session-secret",
				"OAUTH2_PROVIDER_URL":   "http://env-override",
			},
			yamlContent: `
server:
  port: "9090"
oauth2:
  provider_url: "http://yaml-provider"
session:
  cookie_name: "test_cookie"
`,
			wantErr: false,
			check: func(t *testing.T, cfg *Config) {
				if cfg.Server.Port != "9090" {
					t.Errorf("expected port 9090, got %s", cfg.Server.Port)
				}
				if cfg.Session.CookieName != "test_cookie" {
					t.Errorf("expected cookie name test_cookie, got %s", cfg.Session.CookieName)
				}
				if cfg.OAuth2.ProviderURL != "http://env-override" {
					t.Errorf("expected env overridden provider_url, got %s", cfg.OAuth2.ProviderURL)
				}
			},
		},
		{
			name: "Missing Required Missing Session Secret",
			env: map[string]string{
				"OAUTH2_CLIENT_SECRET": "dummy-oauth2-secret",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				t.Setenv(k, v)
			}

			configPath := ""
			if tt.yamlContent != "" {
				tmpDir := t.TempDir()
				configPath = filepath.Join(tmpDir, "config.yaml")
				if err := os.WriteFile(configPath, []byte(tt.yamlContent), 0644); err != nil {
					t.Fatalf("failed to write dummy config: %v", err)
				}
			}

			cfg, err := Load(configPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
