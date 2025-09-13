package nango

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	config := Config{
		APIKey: "test-key",
	}

	client := NewClient(config)

	if client.apiKey != "test-key" {
		t.Errorf("Expected API key 'test-key', got '%s'", client.apiKey)
	}

	if client.baseURL != "https://api.nango.dev" {
		t.Errorf("Expected default base URL 'https://api.nango.dev', got '%s'", client.baseURL)
	}

	if client.httpClient.Timeout != 30*time.Second {
		t.Errorf("Expected default timeout 30s, got %v", client.httpClient.Timeout)
	}
}

func TestNewClientWithCustomConfig(t *testing.T) {
	config := Config{
		BaseURL: "https://custom.api.com",
		APIKey:  "custom-key",
		Timeout: 60 * time.Second,
	}

	client := NewClient(config)

	if client.baseURL != "https://custom.api.com" {
		t.Errorf("Expected base URL 'https://custom.api.com', got '%s'", client.baseURL)
	}

	if client.httpClient.Timeout != 60*time.Second {
		t.Errorf("Expected timeout 60s, got %v", client.httpClient.Timeout)
	}
}

func TestListIntegrations(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/integrations" {
			t.Errorf("Expected path '/integrations', got '%s'", r.URL.Path)
		}

		if r.Header.Get("Authorization") != "Bearer test-key" {
			t.Errorf("Expected Authorization header 'Bearer test-key', got '%s'", r.Header.Get("Authorization"))
		}

		integrations := []Integration{
			{
				ID:       "123",
				Name:     "Test Integration",
				Provider: "github",
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(integrations)
	}))
	defer server.Close()

	config := Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	integrations, err := client.ListIntegrations(context.Background())
	if err != nil {
		t.Fatalf("ListIntegrations failed: %v", err)
	}

	if len(integrations) != 1 {
		t.Errorf("Expected 1 integration, got %d", len(integrations))
	}

	if integrations[0].ID != "123" {
		t.Errorf("Expected integration ID '123', got '%s'", integrations[0].ID)
	}
}

func TestGetIntegration(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/integrations/123" {
			t.Errorf("Expected path '/integrations/123', got '%s'", r.URL.Path)
		}

		integration := Integration{
			ID:       "123",
			Name:     "Test Integration",
			Provider: "github",
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(integration)
	}))
	defer server.Close()

	config := Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	integration, err := client.GetIntegration(context.Background(), "123")
	if err != nil {
		t.Fatalf("GetIntegration failed: %v", err)
	}

	if integration.ID != "123" {
		t.Errorf("Expected integration ID '123', got '%s'", integration.ID)
	}
}

func TestGetIntegrationEmptyID(t *testing.T) {
	config := Config{
		APIKey: "test-key",
	}
	client := NewClient(config)

	_, err := client.GetIntegration(context.Background(), "")
	if err == nil {
		t.Error("Expected error for empty integration ID, got nil")
	}

	if err.Error() != "integration ID cannot be empty" {
		t.Errorf("Expected error 'integration ID cannot be empty', got '%s'", err.Error())
	}
}

func TestGetIntegrationNotFound(t *testing.T) {
	// Mock server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	config := Config{
		BaseURL: server.URL,
		APIKey:  "test-key",
	}
	client := NewClient(config)

	_, err := client.GetIntegration(context.Background(), "nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent integration, got nil")
	}

	expectedError := "integration with ID nonexistent not found"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
