package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_DefaultAndEnv(t *testing.T) {
	// 環境変数をモック
	os.Setenv("OAUTH2_CLIENT_SECRET", "dummy-oauth2-secret")
	os.Setenv("SESSION_COOKIE_SECRET", "dummy-session-secret")
	os.Setenv("SERVER_PORT", "8080")
	defer os.Unsetenv("OAUTH2_CLIENT_SECRET")
	defer os.Unsetenv("SESSION_COOKIE_SECRET")
	defer os.Unsetenv("SERVER_PORT")

	cfg, err := Load("")
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected port 8080, got %s", cfg.Server.Port)
	}
	if cfg.OAuth2.ClientSecret != "dummy-oauth2-secret" {
		t.Errorf("expected secret dummy-oauth2-secret, got %s", cfg.OAuth2.ClientSecret)
	}
}

func TestLoad_YamlAndEnv(t *testing.T) {
	os.Setenv("OAUTH2_CLIENT_SECRET", "dummy-oauth2-secret")
	os.Setenv("SESSION_COOKIE_SECRET", "dummy-session-secret")
	os.Setenv("OAUTH2_PROVIDER_URL", "http://env-override") // 環境変数で上書きされるか確認
	defer os.Unsetenv("OAUTH2_CLIENT_SECRET")
	defer os.Unsetenv("SESSION_COOKIE_SECRET")
	defer os.Unsetenv("OAUTH2_PROVIDER_URL")

	yamlContent := []byte(`
server:
  port: "9090"
oauth2:
  provider_url: "http://yaml-provider"
session:
  cookie_name: "test_cookie"
`)

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")
	if err := os.WriteFile(configPath, yamlContent, 0644); err != nil {
		t.Fatalf("failed to write dummy config: %v", err)
	}

	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// YAMLで設定された値
	if cfg.Server.Port != "9090" {
		t.Errorf("expected port 9090, got %s", cfg.Server.Port)
	}
	if cfg.Session.CookieName != "test_cookie" {
		t.Errorf("expected cookie name test_cookie, got %s", cfg.Session.CookieName)
	}

	// YAMLで設定されたが、環境変数で上書きされた値
	if cfg.OAuth2.ProviderURL != "http://env-override" {
		t.Errorf("expected env overridden provider_url, got %s", cfg.OAuth2.ProviderURL)
	}
}

func TestLoad_MissingRequired(t *testing.T) {
	// シークレットをわざと設定しない
	os.Setenv("OAUTH2_CLIENT_SECRET", "dummy-oauth2-secret")
	defer os.Unsetenv("OAUTH2_CLIENT_SECRET")
	// SESSION_COOKIE_SECRET is missing

	_, err := Load("")
	if err == nil {
		t.Error("expected error due to missing required config, got nil")
	}
}
