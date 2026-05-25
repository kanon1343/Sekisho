package session

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

const CookieName = "_sekisho_session"

var (
	ErrNoSession      = errors.New("no session found")
	ErrSessionExpired = errors.New("session expired")
)

// Session represents the authenticated user data.
type Session struct {
	Subject           string    `json:"sub"`
	Email             string    `json:"email"`
	PreferredUsername string    `json:"preferred_username"`
	Groups            []string  `json:"groups,omitempty"`
	AccessToken       string    `json:"access_token"`
	IDToken           string    `json:"id_token"`
	ExpiresAt         time.Time `json:"expires_at"`
}

// IsExpired checks if the session is expired.
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// Store handles session retrieval and storage in HTTP cookies.
type Store struct {
	cipher *Cipher
	secure bool
}

// NewStore creates a new Store.
func NewStore(cipher *Cipher, secure bool) *Store {
	return &Store{
		cipher: cipher,
		secure: secure,
	}
}

// Get retrieves the session from the request cookie.
func (s *Store) Get(r *http.Request) (*Session, error) {
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return nil, ErrNoSession
		}
		return nil, err
	}

	decoded, err := base64.RawURLEncoding.DecodeString(cookie.Value)
	if err != nil {
		return nil, err
	}

	plaintext, err := s.cipher.Decrypt(decoded)
	if err != nil {
		return nil, err
	}

	var sess Session
	if err := json.Unmarshal(plaintext, &sess); err != nil {
		return nil, err
	}

	if sess.IsExpired() {
		return nil, ErrSessionExpired
	}

	return &sess, nil
}

// Set encrypts and stores the session in a cookie.
func (s *Store) Set(w http.ResponseWriter, sess *Session) error {
	plaintext, err := json.Marshal(sess)
	if err != nil {
		return err
	}

	ciphertext, err := s.cipher.Encrypt(plaintext)
	if err != nil {
		return err
	}

	encoded := base64.RawURLEncoding.EncodeToString(ciphertext)

	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    encoded,
		Path:     "/",
		Expires:  sess.ExpiresAt,
		Secure:   s.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// Clear removes the session cookie.
func (s *Store) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   s.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
