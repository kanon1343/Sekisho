package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v3"
)

// ServerConfig はHTTPサーバーの設定を保持します。
type ServerConfig struct {
	Port string `yaml:"port" envconfig:"SERVER_PORT"`
}

// OAuth2Config はOAuth2プロバイダの設定を保持します。
type OAuth2Config struct {
	ProviderURL  string `yaml:"provider_url" envconfig:"OAUTH2_PROVIDER_URL"`
	ClientID     string `yaml:"client_id" envconfig:"OAUTH2_CLIENT_ID"`
	ClientSecret string `yaml:"client_secret" envconfig:"OAUTH2_CLIENT_SECRET"`
	RedirectURL  string `yaml:"redirect_url" envconfig:"OAUTH2_REDIRECT_URL"`
}

// SessionConfig はセッションとCookieの設定を保持します。
type SessionConfig struct {
	CookieName     string `yaml:"cookie_name" envconfig:"SESSION_COOKIE_NAME"`
	CookieSecret   string `yaml:"cookie_secret" envconfig:"SESSION_COOKIE_SECRET"`
	CookieDomain   string `yaml:"cookie_domain" envconfig:"SESSION_COOKIE_DOMAIN"`
	CookieSecure   bool   `yaml:"cookie_secure" envconfig:"SESSION_COOKIE_SECURE"`
	CookieHTTPOnly bool   `yaml:"cookie_http_only" envconfig:"SESSION_COOKIE_HTTP_ONLY"`
}

// Config はアプリケーション全体の設定を保持します。
type Config struct {
	Server  ServerConfig  `yaml:"server"`
	OAuth2  OAuth2Config  `yaml:"oauth2"`
	Session SessionConfig `yaml:"session"`
}

// DefaultConfig はデフォルト設定を返します。
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "4180",
		},
		OAuth2: OAuth2Config{
			ProviderURL: "http://localhost:8080/realms/sekisho",
			ClientID:    "sekisho-proxy",
			RedirectURL: "https://localhost:4180/oauth2/callback",
		},
		Session: SessionConfig{
			CookieName:     "_sekisho_session",
			CookieDomain:   "localhost",
			CookieSecure:   true,
			CookieHTTPOnly: true,
		},
	}
}

// Load は設定ファイルを読み込み、環境変数で上書きして設定を構築します。
// 設定ファイルパスが空の場合はデフォルト設定を使用し、環境変数でのみ上書きします。
func Load(configPath string) (*Config, error) {
	cfg := DefaultConfig()

	if configPath != "" {
		b, err := os.ReadFile(configPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := yaml.Unmarshal(b, cfg); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
		}
	}

	if err := envconfig.Process("", cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment variables: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return cfg, nil
}

// Validate は設定値の妥当性をチェックします。
func (c *Config) Validate() error {
	if c.OAuth2.ClientSecret == "" {
		return fmt.Errorf("OAUTH2_CLIENT_SECRET is required")
	}
	if c.Session.CookieSecret == "" {
		return fmt.Errorf("SESSION_COOKIE_SECRET is required")
	}
	return nil
}
