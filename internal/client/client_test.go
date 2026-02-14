package client

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("https://opnsense.example.com", "testkey", "testsecret", false)

	if client == nil {
		t.Fatal("Expected client to be non-nil")
	}

	if client.BaseURL != "https://opnsense.example.com" {
		t.Errorf("Expected BaseURL to be 'https://opnsense.example.com', got %s", client.BaseURL)
	}

	if client.APIKey != "testkey" {
		t.Errorf("Expected APIKey to be 'testkey', got %s", client.APIKey)
	}

	if client.APISecret != "testsecret" {
		t.Errorf("Expected APISecret to be 'testsecret', got %s", client.APISecret)
	}
}

func TestDoRequest(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication
		username, password, ok := r.BasicAuth()
		if !ok || username != "testkey" || password != "testsecret" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return success response
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "testkey", "testsecret", true)

	// Test GET request
	resp, err := client.Get("/test")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := `{"result": "success"}`
	if string(resp) != expected {
		t.Errorf("Expected response '%s', got '%s'", expected, string(resp))
	}
}

func TestDoRequestUnauthorized(t *testing.T) {
	// Create a test server that always returns unauthorized
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "unauthorized"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL, "wrongkey", "wrongsecret", true)

	// Test GET request
	_, err := client.Get("/test")
	if err == nil {
		t.Fatal("Expected error for unauthorized request, got nil")
	}
}
