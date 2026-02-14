package resources

import (
	"testing"

	"github.com/yknx4/terraform-provider-opnsense/internal/client"
)

func TestRouteResource_Schema(t *testing.T) {
	r := NewRouteResource()

	if r == nil {
		t.Fatal("Expected resource to be non-nil")
	}

	// Basic test to ensure resource can be instantiated
	// More detailed tests would require a mock OPNsense API
}

func TestRouteResource_Configure(t *testing.T) {
	r := &RouteResource{}

	// Test with valid client
	c := client.NewClient("https://test.example.com", "key", "secret", true)
	r.client = c

	if r.client == nil {
		t.Fatal("Expected client to be set")
	}

	if r.client.BaseURL != "https://test.example.com" {
		t.Errorf("Expected BaseURL to be 'https://test.example.com', got %s", r.client.BaseURL)
	}
}
