// Package nango provides a Go client library for interacting with Nango APIs.
package nango

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client represents a Nango API client.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config holds configuration for the Nango client.
type Config struct {
	BaseURL string
	APIKey  string
	Timeout time.Duration
}

// NewClient creates a new Nango client with the provided configuration.
func NewClient(config Config) *Client {
	if config.BaseURL == "" {
		config.BaseURL = "https://api.nango.dev"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		baseURL: config.BaseURL,
		apiKey:  config.APIKey,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// Integration represents a Nango integration.
type Integration struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Provider  string `json:"provider"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ListIntegrations retrieves all integrations for the authenticated account.
func (c *Client) ListIntegrations(ctx context.Context) ([]Integration, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/integrations", nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var integrations []Integration
	if err := json.NewDecoder(resp.Body).Decode(&integrations); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return integrations, nil
}

// GetIntegration retrieves a specific integration by ID.
func (c *Client) GetIntegration(ctx context.Context, id string) (*Integration, error) {
	if id == "" {
		return nil, fmt.Errorf("integration ID cannot be empty")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/integrations/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("integration with ID %s not found", id)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var integration Integration
	if err := json.NewDecoder(resp.Body).Decode(&integration); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return &integration, nil
}
