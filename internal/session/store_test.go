package session_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"sekisho/internal/session"
)

func TestStore_SetAndGet(t *testing.T) {
	t.Parallel()

	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, err := session.NewCipher(secret)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}
	store := session.NewStore(cipher, false)

	sess := &session.Session{
		Subject:     "user123",
		Email:       "user@example.com",
		AccessToken: "access-token-abc",
		IDToken:     "id-token-xyz",
		ExpiresAt:   time.Now().Add(1 * time.Hour),
	}

	w := httptest.NewRecorder()
	if err := store.Set(w, sess); err != nil {
		t.Fatalf("Store.Set failed: %v", err)
	}

	res := w.Result()
	cookies := res.Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookie to be set")
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range cookies {
		r.AddCookie(c)
	}

	retrieved, err := store.Get(r)
	if err != nil {
		t.Fatalf("Store.Get failed: %v", err)
	}

	if retrieved.Subject != sess.Subject {
		t.Errorf("Expected Subject %s, got %s", sess.Subject, retrieved.Subject)
	}
	if retrieved.Email != sess.Email {
		t.Errorf("Expected Email %s, got %s", sess.Email, retrieved.Email)
	}
}

func TestStore_GetNoSession(t *testing.T) {
	t.Parallel()

	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, _ := session.NewCipher(secret)
	store := session.NewStore(cipher, false)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	_, err := store.Get(r)
	if err != session.ErrNoSession {
		t.Errorf("Expected ErrNoSession, got %v", err)
	}
}

func TestStore_GetExpired(t *testing.T) {
	t.Parallel()

	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, _ := session.NewCipher(secret)
	store := session.NewStore(cipher, false)

	sess := &session.Session{
		Subject:   "user123",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}

	w := httptest.NewRecorder()
	if err := store.Set(w, sess); err != nil {
		t.Fatalf("Store.Set failed: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	for _, c := range w.Result().Cookies() {
		r.AddCookie(c)
	}

	_, err := store.Get(r)
	if err != session.ErrSessionExpired {
		t.Errorf("Expected ErrSessionExpired, got %v", err)
	}
}

func TestStore_Clear(t *testing.T) {
	t.Parallel()

	secret := []byte("super-secret-key-32-bytes-long!")
	cipher, _ := session.NewCipher(secret)
	store := session.NewStore(cipher, false)

	w := httptest.NewRecorder()
	store.Clear(w)

	cookies := w.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("Expected cookie to be set")
	}

	c := cookies[0]
	if c.Name != session.CookieName {
		t.Errorf("Expected cookie name %s, got %s", session.CookieName, c.Name)
	}
	if c.Value != "" {
		t.Errorf("Expected empty cookie value, got %s", c.Value)
	}
	if c.MaxAge != -1 {
		t.Errorf("Expected MaxAge -1, got %d", c.MaxAge)
	}
}

