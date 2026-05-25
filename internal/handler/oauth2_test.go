package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerateState(t *testing.T) {
	t.Parallel()
	state1, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState failed: %v", err)
	}
	state2, err := GenerateState()
	if err != nil {
		t.Fatalf("GenerateState failed: %v", err)
	}
	if state1 == state2 {
		t.Errorf("GenerateState returned same state twice: %s", state1)
	}
	if len(state1) < 32 {
		t.Errorf("GenerateState length is too short: got %d", len(state1))
	}
}

func TestValidateState(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		queryState string
		cookieVal  string
		wantErr    bool
	}{
		{
			name:       "valid state",
			queryState: "test-state",
			cookieVal:  "test-state",
			wantErr:    false,
		},
		{
			name:       "missing query state",
			queryState: "",
			cookieVal:  "test-state",
			wantErr:    true,
		},
		{
			name:       "missing cookie state",
			queryState: "test-state",
			cookieVal:  "",
			wantErr:    true,
		},
		{
			name:       "state mismatch",
			queryState: "test-state",
			cookieVal:  "wrong-state",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/callback?state="+tt.queryState, nil)
			if tt.cookieVal != "" {
				req.AddCookie(&http.Cookie{
					Name:  StateCookieName,
					Value: tt.cookieVal,
				})
			}

			err := ValidateState(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateRedirectURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		rd          string
		allowedHost string
		want        bool
	}{
		{
			name:        "valid absolute URL",
			rd:          "https://example.com/app",
			allowedHost: "example.com",
			want:        true,
		},
		{
			name:        "invalid absolute URL host",
			rd:          "https://evil.com/app",
			allowedHost: "example.com",
			want:        false,
		},
		{
			name:        "valid relative path",
			rd:          "/dashboard",
			allowedHost: "example.com",
			want:        true,
		},
		{
			name:        "invalid protocol relative URL",
			rd:          "//evil.com",
			allowedHost: "example.com",
			want:        false,
		},
		{
			name:        "invalid empty URL",
			rd:          "",
			allowedHost: "example.com",
			want:        false,
		},
		{
			name:        "invalid relative path without slash",
			rd:          "dashboard",
			allowedHost: "example.com",
			want:        false,
		},
		{
			name:        "valid root path",
			rd:          "/",
			allowedHost: "example.com",
			want:        true,
		},
		{
			name:        "javascript scheme",
			rd:          "javascript:alert(1)",
			allowedHost: "example.com",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateRedirectURL(tt.rd, tt.allowedHost); got != tt.want {
				t.Errorf("ValidateRedirectURL(%q, %q) = %v, want %v", tt.rd, tt.allowedHost, got, tt.want)
			}
		})
	}
}
