package oidc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// OIDCClient はOIDCプロバイダとのやり取りを行うクライアントのインターフェースです。
type OIDCClient interface {
	GetProviderMetadata(ctx context.Context) (*ProviderMetadata, error)
	// TODO: AuthURLの生成やTokenの取得メソッドを今後追加
}

// ProviderMetadata はOIDCのDiscovery Endpointから取得するメタデータです。
type ProviderMetadata struct {
	Issuer                           string   `json:"issuer"`
	AuthorizationEndpoint            string   `json:"authorization_endpoint"`
	TokenEndpoint                    string   `json:"token_endpoint"`
	JwksURI                          string   `json:"jwks_uri"`
	UserInfoEndpoint                 string   `json:"userinfo_endpoint"`
	EndSessionEndpoint               string   `json:"end_session_endpoint"`
	ScopesSupported                  []string `json:"scopes_supported"`
	ResponseTypesSupported           []string `json:"response_types_supported"`
	IdTokenSigningAlgValuesSupported []string `json:"id_token_signing_alg_values_supported"`
}

// client はOIDCClientインターフェースの実装です。
type client struct {
	providerURL string
	httpClient  *http.Client
}

// NewClient は新しいOIDCClientを生成します。
func NewClient(providerURL string, httpClient *http.Client) OIDCClient {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	return &client{
		providerURL: strings.TrimSuffix(providerURL, "/"),
		httpClient:  httpClient,
	}
}

// GetProviderMetadata はDiscovery Endpointからプロバイダのメタデータを取得します。
func (c *client) GetProviderMetadata(ctx context.Context) (*ProviderMetadata, error) {
	discoveryURL := fmt.Sprintf("%s/.well-known/openid-configuration", c.providerURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, discoveryURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for discovery: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute discovery request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from discovery endpoint: %d", resp.StatusCode)
	}

	var metadata ProviderMetadata
	if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
		return nil, fmt.Errorf("failed to decode discovery response: %w", err)
	}

	return &metadata, nil
}
