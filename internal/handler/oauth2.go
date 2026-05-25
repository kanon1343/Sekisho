package handler

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"sekisho/internal/config"
	"sekisho/internal/oidc"
	"sekisho/internal/session"
)

// OAuth2Handler はOAuth2関連のHTTPエンドポイントを処理します。
type OAuth2Handler struct {
	Config *config.OAuth2Config
	OIDC   oidc.OIDCClient
	Cipher *session.Cipher
}

// NewOAuth2Handler は新しいOAuth2Handlerを作成します。
func NewOAuth2Handler(cfg *config.OAuth2Config, oidcClient oidc.OIDCClient, cipher *session.Cipher) *OAuth2Handler {
	return &OAuth2Handler{
		Config: cfg,
		OIDC:   oidcClient,
		Cipher: cipher,
	}
}

// generateRandomString は暗号学的に安全なランダムな文字列を生成します。
func generateRandomString(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// HandleStart は /oauth2/start エンドポイントを処理します。
// 認可URLを生成し、CSRF情報をCookieに保存してIdPにリダイレクトします。
func (h *OAuth2Handler) HandleStart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metadata, err := h.OIDC.GetProviderMetadata(ctx)
	if err != nil {
		http.Error(w, "Failed to get OIDC provider metadata", http.StatusInternalServerError)
		return
	}

	state, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate state", http.StatusInternalServerError)
		return
	}

	nonce, err := generateRandomString(32)
	if err != nil {
		http.Error(w, "Failed to generate nonce", http.StatusInternalServerError)
		return
	}

	codeVerifier, err := oidc.GenerateCodeVerifier()
	if err != nil {
		http.Error(w, "Failed to generate code verifier", http.StatusInternalServerError)
		return
	}
	codeChallenge := oidc.GenerateCodeChallenge(codeVerifier)

	rd := r.URL.Query().Get("rd")
	if rd == "" {
		rd = "/"
	}

	cookieData := fmt.Sprintf("%s|%s|%s|%s", state, nonce, codeVerifier, rd)

	encryptedData, err := h.Cipher.Encrypt([]byte(cookieData))
	if err != nil {
		http.Error(w, "Failed to encrypt CSRF cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "_sekisho_csrf",
		Value:    base64.RawURLEncoding.EncodeToString(encryptedData),
		MaxAge:   300,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})

	authURL, err := url.Parse(metadata.AuthorizationEndpoint)
	if err != nil {
		http.Error(w, "Failed to parse authorization endpoint", http.StatusInternalServerError)
		return
	}

	q := authURL.Query()
	q.Set("client_id", h.Config.ClientID)
	q.Set("redirect_uri", h.Config.RedirectURL)
	q.Set("response_type", "code")
	q.Set("scope", "openid profile email")
	q.Set("state", state)
	q.Set("nonce", nonce)
	q.Set("code_challenge", codeChallenge)
	q.Set("code_challenge_method", "S256")
	authURL.RawQuery = q.Encode()

	http.Redirect(w, r, authURL.String(), http.StatusFound)
}
