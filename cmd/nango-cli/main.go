// nango-cli is a command-line interface for interacting with Nango APIs.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/saltmueller/go-nango/internal/config"
	"github.com/saltmueller/go-nango/pkg/nango"
)

func main() {
	var (
		command       = flag.String("command", "list", "Command to execute (list, get)")
		integrationID = flag.String("id", "", "Integration ID (required for 'get' command)")
		apiKey        = flag.String("api-key", "", "Nango API key (overrides NANGO_API_KEY env var)")
		baseURL       = flag.String("base-url", "", "Nango base URL (overrides NANGO_BASE_URL env var)")
		timeout       = flag.Duration("timeout", 30*time.Second, "Request timeout")
		help          = flag.Bool("help", false, "Show help message")
	)
	flag.Parse()

	if *help {
		showHelp()
		return
	}

	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Printf("Warning: Failed to load config from environment: %v", err)
		// Create a minimal config if env loading fails
		cfg = &config.AppConfig{
			NangoAPIKey:  "",
			NangoBaseURL: "https://api.nango.dev",
			Timeout:      30 * time.Second,
		}
	}

	// Override with command-line flags if provided
	if *apiKey != "" {
		cfg.NangoAPIKey = *apiKey
	}
	if *baseURL != "" {
		cfg.NangoBaseURL = *baseURL
	}
	if *timeout != 30*time.Second {
		cfg.Timeout = *timeout
	}

	// Check for required API key
	if cfg.NangoAPIKey == "" {
		fmt.Fprintf(os.Stderr, "Error: API key is required. Set NANGO_API_KEY environment variable or use -api-key flag.\n")
		os.Exit(1)
	}

	// Create Nango client
	nangoConfig := nango.Config{
		BaseURL: cfg.NangoBaseURL,
		APIKey:  cfg.NangoAPIKey,
		Timeout: cfg.Timeout,
	}
	client := nango.NewClient(nangoConfig)

	ctx := context.Background()

	switch *command {
	case "list":
		if err := listIntegrations(ctx, client); err != nil {
			fmt.Fprintf(os.Stderr, "Error listing integrations: %v\n", err)
			os.Exit(1)
		}
	case "get":
		if *integrationID == "" {
			fmt.Fprintf(os.Stderr, "Error: Integration ID is required for 'get' command. Use -id flag.\n")
			os.Exit(1)
		}
		if err := getIntegration(ctx, client, *integrationID); err != nil {
			fmt.Fprintf(os.Stderr, "Error getting integration: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'. Use -help for available commands.\n", *command)
		os.Exit(1)
	}
}

func listIntegrations(ctx context.Context, client *nango.Client) error {
	integrations, err := client.ListIntegrations(ctx)
	if err != nil {
		return err
	}

	if len(integrations) == 0 {
		fmt.Println("No integrations found.")
		return nil
	}

	fmt.Printf("Found %d integration(s):\n\n", len(integrations))
	for _, integration := range integrations {
		output, err := json.MarshalIndent(integration, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling integration: %w", err)
		}
		fmt.Println(string(output))
		fmt.Println()
	}

	return nil
}

func getIntegration(ctx context.Context, client *nango.Client, id string) error {
	integration, err := client.GetIntegration(ctx, id)
	if err != nil {
		return err
	}

	output, err := json.MarshalIndent(integration, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling integration: %w", err)
	}

	fmt.Println(string(output))
	return nil
}

func showHelp() {
	fmt.Printf(`nango-cli - Command-line interface for Nango APIs

Usage:
  nango-cli [options] -command <command>

Commands:
  list    List all integrations
  get     Get a specific integration by ID (requires -id flag)

Options:
  -command string
        Command to execute (default "list")
  -id string
        Integration ID (required for 'get' command)
  -api-key string
        Nango API key (overrides NANGO_API_KEY env var)
  -base-url string
        Nango base URL (overrides NANGO_BASE_URL env var)
  -timeout duration
        Request timeout (default 30s)
  -help
        Show this help message

Environment Variables:
  NANGO_API_KEY    Your Nango API key (required)
  NANGO_BASE_URL   Nango API base URL (default: https://api.nango.dev)

Examples:
  # List all integrations
  nango-cli -command list

  # Get a specific integration
  nango-cli -command get -id 123

  # Use custom API key and base URL
  nango-cli -api-key your-key -base-url https://custom.api.com -command list
`)
}
