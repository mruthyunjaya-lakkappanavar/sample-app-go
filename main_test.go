package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestGreetWithName(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/api/greet?name=Alice", nil)
	if err != nil {
		t.Fatal(err)
	}

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

	req, err := http.NewRequest("GET", "/api/greet", nil)
	if err != nil {
		t.Fatal(err)
	}

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

func TestStatusEndpoint(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}

	var resp StatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.Version == "" {
		t.Error("expected version to be non-empty")
	}

	if resp.GoVersion == "" {
		t.Error("expected go_version to be non-empty")
	}

	if resp.Uptime == "" {
		t.Error("expected uptime to be non-empty")
	}
}

func TestStatusEndpointContentType(t *testing.T) {
	router := NewRouter()

	req, err := http.NewRequest("GET", "/api/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestRequestCounter(t *testing.T) {
	router := NewRouter()

	// Make a few requests to increment counter
	for i := 0; i < 3; i++ {
		req, _ := http.NewRequest("GET", "/api/status", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
	}

	// Verify counter is at least 3
	req, _ := http.NewRequest("GET", "/api/status", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	var resp StatusResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp.RequestCount < 3 {
		t.Errorf("expected request_count >= 3, got %d", resp.RequestCount)
	}
}
