package handler

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// StateCookieName is the name of the cookie used to store the OAuth2 state parameter.
	StateCookieName = "_sekisho_oauth_state"
	// StateCookieTTL is the time-to-live for the state cookie.
	StateCookieTTL = 10 * time.Minute
)

// GenerateState generates a secure random state string for OAuth2 CSRF protection.
func GenerateState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// SetStateCookie sets the state string into a secure, HttpOnly, SameSite=Lax cookie.
func SetStateCookie(w http.ResponseWriter, state string, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     StateCookieName,
		Value:    state,
		Path:     "/",
		Expires:  time.Now().Add(StateCookieTTL),
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ClearStateCookie removes the state cookie.
func ClearStateCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     StateCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}

// ValidateState checks if the state from the request query matches the state from the cookie.
func ValidateState(r *http.Request) error {
	queryState := r.URL.Query().Get("state")
	if queryState == "" {
		return errors.New("missing state parameter in request")
	}

	cookie, err := r.Cookie(StateCookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return errors.New("missing state cookie")
		}
		return err
	}

	if cookie.Value != queryState {
		return errors.New("state mismatch")
	}

	return nil
}

// ValidateRedirectURL checks if the provided URL is a safe redirect destination
// to prevent Open Redirect vulnerabilities. It allows relative URLs (starting with /)
// or absolute URLs that match the allowed host.
func ValidateRedirectURL(rd string, allowedHost string) bool {
	if rd == "" {
		return false
	}

	u, err := url.Parse(rd)
	if err != nil {
		return false
	}

	// If it's a relative URL (no host)
	if u.Host == "" {
		// It must start with a single '/' to be a valid relative path on the same domain.
		// It prevents protocol-relative URLs like "//evil.com" which start with "//".
		if !strings.HasPrefix(rd, "/") || strings.HasPrefix(rd, "//") {
			return false
		}
		return true
	}

	// If it's an absolute URL, the host must match the allowed host exactly.
	if u.Host == allowedHost {
		return true
	}

	return false
}
