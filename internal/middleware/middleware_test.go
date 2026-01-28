package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewResponseWriter(t *testing.T) {
	tests := []struct {
		name              string
		wantDefaultStatus int
		wantWroteHeader   bool
	}{
		{
			name:              "default status is 200 and wroteHeader is false",
			wantDefaultStatus: http.StatusOK,
			wantWroteHeader:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rw := newResponseWriter(rec)

			if rw.statusCode != tt.wantDefaultStatus {
				t.Errorf("statusCode = %d, want %d", rw.statusCode, tt.wantDefaultStatus)
			}

			if rw.wroteHeader != tt.wantWroteHeader {
				t.Errorf("wroteHeader = %v, want %v", rw.wroteHeader, tt.wantWroteHeader)
			}

			if rw.ResponseWriter != rec {
				t.Error("ResponseWriter not set correctly")
			}
		})
	}
}

func TestResponseWriter_WriteHeader(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{
			name:       "200 OK",
			statusCode: http.StatusOK,
		},
		{
			name:       "201 Created",
			statusCode: http.StatusCreated,
		},
		{
			name:       "400 Bad Request",
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "404 Not Found",
			statusCode: http.StatusNotFound,
		},
		{
			name:       "500 Internal Server Error",
			statusCode: http.StatusInternalServerError,
		},
		{
			name:       "503 Service Unavailable",
			statusCode: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			rw := newResponseWriter(rec)

			rw.WriteHeader(tt.statusCode)

			if rw.statusCode != tt.statusCode {
				t.Errorf("statusCode = %d, want %d", rw.statusCode, tt.statusCode)
			}

			if !rw.wroteHeader {
				t.Error("wroteHeader should be true after WriteHeader")
			}

			if rec.Code != tt.statusCode {
				t.Errorf("recorder Code = %d, want %d", rec.Code, tt.statusCode)
			}
		})
	}
}

func TestResponseWriter_WriteHeader_OnlyFirstCall(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := newResponseWriter(rec)

	// First call should set status code
	rw.WriteHeader(http.StatusCreated)
	if rw.statusCode != http.StatusCreated {
		t.Errorf("statusCode = %d, want %d", rw.statusCode, http.StatusCreated)
	}

	// Second call should not update the captured status code
	rw.WriteHeader(http.StatusInternalServerError)
	if rw.statusCode != http.StatusCreated {
		t.Errorf("statusCode changed to %d, want %d (should not change)", rw.statusCode, http.StatusCreated)
	}
}

func TestMetricsMiddleware(t *testing.T) {
	tests := []struct {
		name          string
		method        string
		path          string
		handlerStatus int
		handlerBody   string
		wantStatus    int
		wantBody      string
	}{
		{
			name:          "GET request returns 200",
			method:        http.MethodGet,
			path:          "/test",
			handlerStatus: http.StatusOK,
			handlerBody:   "OK",
			wantStatus:    http.StatusOK,
			wantBody:      "OK",
		},
		{
			name:          "POST request returns 201",
			method:        http.MethodPost,
			path:          "/api/resource",
			handlerStatus: http.StatusCreated,
			handlerBody:   `{"id":1}`,
			wantStatus:    http.StatusCreated,
			wantBody:      `{"id":1}`,
		},
		{
			name:          "GET request returns 404",
			method:        http.MethodGet,
			path:          "/notfound",
			handlerStatus: http.StatusNotFound,
			handlerBody:   "Not Found",
			wantStatus:    http.StatusNotFound,
			wantBody:      "Not Found",
		},
		{
			name:          "DELETE request returns 204",
			method:        http.MethodDelete,
			path:          "/api/resource/1",
			handlerStatus: http.StatusNoContent,
			handlerBody:   "",
			wantStatus:    http.StatusNoContent,
			wantBody:      "",
		},
		{
			name:          "PUT request returns 500",
			method:        http.MethodPut,
			path:          "/api/error",
			handlerStatus: http.StatusInternalServerError,
			handlerBody:   "Internal Error",
			wantStatus:    http.StatusInternalServerError,
			wantBody:      "Internal Error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.handlerStatus)
				if tt.handlerBody != "" {
					_, _ = w.Write([]byte(tt.handlerBody))
				}
			})

			wrapped := MetricsMiddleware(handler)

			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			wrapped.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}

			if rec.Body.String() != tt.wantBody {
				t.Errorf("body = %q, want %q", rec.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestMetricsMiddleware_DefaultStatus(t *testing.T) {
	// Test handler that writes body without explicit WriteHeader
	// This should default to 200 OK
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("implicit 200"))
	})

	wrapped := MetricsMiddleware(handler)

	req := httptest.NewRequest(http.MethodGet, "/implicit", nil)
	rec := httptest.NewRecorder()

	wrapped.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
