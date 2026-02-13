package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// failWriter is a ResponseWriter that always fails on Write,
// used to exercise error-handling branches in handlers.
type failWriter struct{}

func (fw failWriter) Header() http.Header        { return http.Header{} }
func (fw failWriter) Write([]byte) (int, error)   { return 0, fmt.Errorf("write error") }
func (fw failWriter) WriteHeader(int)             {}

// ── Health endpoint ─────────────────────────────────────────

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp HealthResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Status != "ok" {
		t.Errorf("expected status 'ok', got '%s'", resp.Status)
	}

	if resp.Version == "" {
		t.Error("expected version to be non-empty")
	}
}

func TestHealthEndpointContentType(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestHealthHandlerWriteError(t *testing.T) {
	req := httptest.NewRequest("GET", "/health", nil)
	healthHandler(failWriter{}, req)
	// exercises the http.Error branch when json.Encode fails
}

// ── Greet endpoint ──────────────────────────────────────────

func TestGreetWithName(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/api/greet?name=Alice", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp GreetResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := "Hello, Alice!"
	if resp.Message != expected {
		t.Errorf("expected message '%s', got '%s'", expected, resp.Message)
	}
}

func TestGreetWithoutName(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/api/greet", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp GreetResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	expected := "Hello, World!"
	if resp.Message != expected {
		t.Errorf("expected message '%s', got '%s'", expected, resp.Message)
	}
}

func TestGreetHandlerWriteError(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/greet?name=Test", nil)
	greetHandler(failWriter{}, req)
}

// ── Info endpoint ───────────────────────────────────────────

func TestInfoEndpoint(t *testing.T) {
	router := NewRouter()

	req := httptest.NewRequest("GET", "/api/info", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp InfoResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.App != "sample-app-go" {
		t.Errorf("expected app 'sample-app-go', got '%s'", resp.App)
	}

	if resp.GoVersion == "" {
		t.Error("expected GoVersion to be non-empty")
	}
}

func TestInfoHandlerWriteError(t *testing.T) {
	req := httptest.NewRequest("GET", "/api/info", nil)
	infoHandler(failWriter{}, req)
}

// ── getEnv ──────────────────────────────────────────────────

func TestGetEnvWithValue(t *testing.T) {
	t.Setenv("TEST_KEY_DEMO", "custom_value")

	result := getEnv("TEST_KEY_DEMO", "fallback")
	if result != "custom_value" {
		t.Errorf("expected 'custom_value', got '%s'", result)
	}
}

func TestGetEnvWithFallback(t *testing.T) {
	os.Unsetenv("NONEXISTENT_KEY_DEMO")

	result := getEnv("NONEXISTENT_KEY_DEMO", "default")
	if result != "default" {
		t.Errorf("expected 'default', got '%s'", result)
	}
}

// ── newServer ───────────────────────────────────────────────

func TestNewServer(t *testing.T) {
	srv := newServer()

	if srv.Addr != "127.0.0.1:8080" {
		t.Errorf("expected addr '127.0.0.1:8080', got '%s'", srv.Addr)
	}
	if srv.ReadTimeout != 15*time.Second {
		t.Errorf("expected ReadTimeout 15s, got %v", srv.ReadTimeout)
	}
	if srv.WriteTimeout != 15*time.Second {
		t.Errorf("expected WriteTimeout 15s, got %v", srv.WriteTimeout)
	}
	if srv.IdleTimeout != 60*time.Second {
		t.Errorf("expected IdleTimeout 60s, got %v", srv.IdleTimeout)
	}
	if srv.Handler == nil {
		t.Error("expected non-nil handler")
	}
}

func TestNewServerCustomPort(t *testing.T) {
	t.Setenv("PORT", "9090")

	srv := newServer()
	if srv.Addr != "127.0.0.1:9090" {
		t.Errorf("expected addr '127.0.0.1:9090', got '%s'", srv.Addr)
	}
}
