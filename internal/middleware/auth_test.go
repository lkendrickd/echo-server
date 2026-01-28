package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// mockValidator implements APIKeyValidator for testing
type mockValidator struct {
	validKeys map[string]struct{}
}

func newMockValidator(keys ...string) *mockValidator {
	v := &mockValidator{validKeys: make(map[string]struct{})}
	for _, k := range keys {
		v.validKeys[k] = struct{}{}
	}
	return v
}

func (m *mockValidator) ValidateAPIKey(key string) bool {
	_, exists := m.validKeys[key]
	return exists
}

func TestAuthMiddleware(t *testing.T) {
	protectedPrefixes := []string{"/api/"}

	tests := []struct {
		name           string
		path           string
		apiKey         string
		validKeys      []string
		wantStatus     int
		wantError      string
		shouldCallNext bool
	}{
		{
			name:           "valid key on protected path",
			path:           "/api/v1/echo",
			apiKey:         "valid-key",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "missing key on protected path",
			path:           "/api/v1/echo",
			apiKey:         "",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "missing API key",
			shouldCallNext: false,
		},
		{
			name:           "invalid key on protected path",
			path:           "/api/v1/echo",
			apiKey:         "wrong-key",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "invalid API key",
			shouldCallNext: false,
		},
		{
			name:           "unprotected path without key",
			path:           "/health",
			apiKey:         "",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "unprotected path with invalid key still passes",
			path:           "/metrics",
			apiKey:         "wrong-key",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "multiple valid keys - first key",
			path:           "/api/v1/echo",
			apiKey:         "key1",
			validKeys:      []string{"key1", "key2", "key3"},
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "multiple valid keys - middle key",
			path:           "/api/v1/echo",
			apiKey:         "key2",
			validKeys:      []string{"key1", "key2", "key3"},
			wantStatus:     http.StatusOK,
			shouldCallNext: true,
		},
		{
			name:           "case sensitive key validation",
			path:           "/api/v1/echo",
			apiKey:         "VALID-KEY",
			validKeys:      []string{"valid-key"},
			wantStatus:     http.StatusUnauthorized,
			wantError:      "invalid API key",
			shouldCallNext: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextCalled := false
			nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				nextCalled = true
				w.WriteHeader(http.StatusOK)
			})

			validator := newMockValidator(tt.validKeys...)
			middleware := AuthMiddleware(validator, protectedPrefixes)
			handler := middleware(nextHandler)

			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.apiKey != "" {
				req.Header.Set("X-API-Key", tt.apiKey)
			}
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if nextCalled != tt.shouldCallNext {
				t.Errorf("next handler called = %v, want %v", nextCalled, tt.shouldCallNext)
			}

			if tt.wantError != "" {
				var errResp authErrorResponse
				if err := json.NewDecoder(rec.Body).Decode(&errResp); err != nil {
					t.Fatalf("failed to decode error response: %v", err)
				}
				if errResp.Error != tt.wantError {
					t.Errorf("error = %q, want %q", errResp.Error, tt.wantError)
				}

				if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
					t.Errorf("Content-Type = %q, want application/json", ct)
				}
			}
		})
	}
}

func TestIsProtectedPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		prefixes []string
		want     bool
	}{
		{
			name:     "matches first prefix",
			path:     "/api/v1/echo",
			prefixes: []string{"/api/", "/admin/"},
			want:     true,
		},
		{
			name:     "matches second prefix",
			path:     "/admin/users",
			prefixes: []string{"/api/", "/admin/"},
			want:     true,
		},
		{
			name:     "no match",
			path:     "/health",
			prefixes: []string{"/api/", "/admin/"},
			want:     false,
		},
		{
			name:     "empty prefixes",
			path:     "/api/v1/echo",
			prefixes: []string{},
			want:     false,
		},
		{
			name:     "exact prefix match",
			path:     "/api/",
			prefixes: []string{"/api/"},
			want:     true,
		},
		{
			name:     "partial path doesn't match",
			path:     "/ap",
			prefixes: []string{"/api/"},
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isProtectedPath(tt.path, tt.prefixes); got != tt.want {
				t.Errorf("isProtectedPath(%q, %v) = %v, want %v", tt.path, tt.prefixes, got, tt.want)
			}
		})
	}
}

func TestSecureCompare(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{
			name: "equal strings",
			a:    "secret-key-123",
			b:    "secret-key-123",
			want: true,
		},
		{
			name: "different strings",
			a:    "secret-key-123",
			b:    "secret-key-456",
			want: false,
		},
		{
			name: "different lengths",
			a:    "short",
			b:    "longer-string",
			want: false,
		},
		{
			name: "empty strings",
			a:    "",
			b:    "",
			want: true,
		},
		{
			name: "one empty",
			a:    "key",
			b:    "",
			want: false,
		},
		{
			name: "case sensitive",
			a:    "Key",
			b:    "key",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SecureCompare(tt.a, tt.b); got != tt.want {
				t.Errorf("SecureCompare(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
