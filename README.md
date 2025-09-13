# go-nango

A Go library and command-line tools for interacting with Nango APIs.

## Project Structure

```
go-nango/
├── cmd/                    # Executable commands
│   ├── nango-cli/         # Command-line interface
│   └── nango-server/      # HTTP server
├── pkg/                   # Public library packages
│   └── nango/            # Main Nango client library
├── internal/             # Private application packages
│   └── config/          # Configuration utilities
├── examples/            # Usage examples
│   └── basic/          # Basic usage example
├── docs/               # Documentation
├── Makefile           # Build automation
├── go.mod             # Go module file
├── go.sum             # Go module checksums
├── LICENSE            # MIT License
└── README.md          # This file
```

## Features

- **Go Library**: Clean, well-tested client library for Nango APIs
- **CLI Tool**: Command-line interface for managing integrations
- **HTTP Server**: REST API server with Nango integration endpoints
- **Configuration Management**: Environment-based configuration
- **Comprehensive Testing**: Unit tests with mocking
- **Examples**: Ready-to-run usage examples

## Installation

### Prerequisites

- Go 1.21 or later
- Nango API key

### Using Go Install

```bash
# Install the CLI tool
go install github.com/saltmueller/go-nango/cmd/nango-cli@latest

# Install the server
go install github.com/saltmueller/go-nango/cmd/nango-server@latest
```

### Building from Source

```bash
# Clone the repository
git clone https://github.com/saltmueller/go-nango.git
cd go-nango

# Build all components
make build

# Or build specific components
make build-cli
make build-server
```

## Usage

### Library Usage

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/saltmueller/go-nango/pkg/nango"
)

func main() {
    config := nango.Config{
        BaseURL: "https://api.nango.dev",
        APIKey:  "your-api-key",
        Timeout: 30 * time.Second,
    }

    client := nango.NewClient(config)
    
    // List integrations
    integrations, err := client.ListIntegrations(context.Background())
    if err != nil {
        panic(err)
    }

    fmt.Printf("Found %d integrations\n", len(integrations))
}
```

### CLI Tool

```bash
# Set your API key
export NANGO_API_KEY=your-api-key

# List all integrations
nango-cli -command list

# Get a specific integration
nango-cli -command get -id 123

# Use custom base URL
nango-cli -base-url https://custom.api.com -command list

# Show help
nango-cli -help
```

### HTTP Server

```bash
# Set required environment variables
export NANGO_API_KEY=your-api-key
export PORT=8080

# Start the server
nango-server
```

The server provides the following endpoints:

- `GET /health` - Health check
- `GET /integrations` - List all integrations
- `GET /integrations/{id}` - Get specific integration

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `NANGO_API_KEY` | Your Nango API key | - | Yes |
| `NANGO_BASE_URL` | Nango API base URL | `https://api.nango.dev` | No |
| `PORT` | Server port | `8080` | No |
| `TIMEOUT` | Request timeout | `30s` | No |
| `LOG_LEVEL` | Log level | `info` | No |

## Development

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage
```

### Building

```bash
# Build all components
make build

# Clean build artifacts
make clean

# Format code
make fmt

# Run linter
make lint
```

### Running Examples

```bash
# Run basic example
make example
```

## API Documentation

### Client

#### `NewClient(config Config) *Client`

Creates a new Nango client with the provided configuration.

#### `ListIntegrations(ctx context.Context) ([]Integration, error)`

Retrieves all integrations for the authenticated account.

#### `GetIntegration(ctx context.Context, id string) (*Integration, error)`

Retrieves a specific integration by ID.

### Types

```go
type Config struct {
    BaseURL string        // Nango API base URL
    APIKey  string        // API key for authentication
    Timeout time.Duration // Request timeout
}

type Integration struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Provider  string `json:"provider"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run `make lint` and `make test`
6. Submit a pull request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
