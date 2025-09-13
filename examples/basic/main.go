// Example demonstrates basic usage of the go-nango library.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/saltmueller/go-nango/pkg/nango"
)

func main() {
	// Create a new Nango client
	config := nango.Config{
		BaseURL: "https://api.nango.dev",
		APIKey:  "your-api-key-here", // Replace with your actual API key
		Timeout: 30 * time.Second,
	}

	client := nango.NewClient(config)

	ctx := context.Background()

	// Example 1: List all integrations
	fmt.Println("=== Listing Integrations ===")
	integrations, err := client.ListIntegrations(ctx)
	if err != nil {
		log.Printf("Error listing integrations: %v", err)
	} else {
		fmt.Printf("Found %d integrations:\n", len(integrations))
		for _, integration := range integrations {
			fmt.Printf("- ID: %s, Name: %s, Provider: %s\n",
				integration.ID, integration.Name, integration.Provider)
		}
	}

	// Example 2: Get a specific integration (if any exist)
	if len(integrations) > 0 {
		fmt.Println("\n=== Getting Specific Integration ===")
		firstIntegration := integrations[0]
		integration, err := client.GetIntegration(ctx, firstIntegration.ID)
		if err != nil {
			log.Printf("Error getting integration: %v", err)
		} else {
			fmt.Printf("Integration Details:\n")
			fmt.Printf("  ID: %s\n", integration.ID)
			fmt.Printf("  Name: %s\n", integration.Name)
			fmt.Printf("  Provider: %s\n", integration.Provider)
			fmt.Printf("  Created: %s\n", integration.CreatedAt)
			fmt.Printf("  Updated: %s\n", integration.UpdatedAt)
		}
	}

	// Example 3: Error handling for non-existent integration
	fmt.Println("\n=== Error Handling Example ===")
	_, err = client.GetIntegration(ctx, "non-existent-id")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
	}
}
