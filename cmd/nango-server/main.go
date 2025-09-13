// nango-server is a simple HTTP server that provides REST endpoints for Nango APIs.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/saltmueller/go-nango/internal/config"
	"github.com/saltmueller/go-nango/pkg/nango"
)

type server struct {
	nangoClient *nango.Client
}

func main() {
	// Load configuration
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// Create Nango client
	nangoConfig := nango.Config{
		BaseURL: cfg.NangoBaseURL,
		APIKey:  cfg.NangoAPIKey,
		Timeout: cfg.Timeout,
	}
	nangoClient := nango.NewClient(nangoConfig)

	// Create server
	srv := &server{
		nangoClient: nangoClient,
	}

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/health", srv.healthHandler)
	mux.HandleFunc("/integrations", srv.integrationsHandler)
	mux.HandleFunc("/integrations/", srv.integrationHandler)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func (s *server) healthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	response := map[string]string{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *server) integrationsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	integrations, err := s.nangoClient.ListIntegrations(r.Context())
	if err != nil {
		log.Printf("Failed to list integrations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(integrations)
}

func (s *server) integrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract integration ID from path
	path := strings.TrimPrefix(r.URL.Path, "/integrations/")
	if path == "" {
		http.Error(w, "Integration ID is required", http.StatusBadRequest)
		return
	}

	integration, err := s.nangoClient.GetIntegration(r.Context(), path)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Integration not found", http.StatusNotFound)
			return
		}
		log.Printf("Failed to get integration %s: %v", path, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(integration)
}
