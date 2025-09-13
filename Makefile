.PHONY: all build test clean install lint fmt vet

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOVET=$(GOCMD) vet

# Binary names
CLI_BINARY_NAME=nango-cli
SERVER_BINARY_NAME=nango-server
BINARY_DIR=bin

# Build targets
all: test build

build: build-cli build-server

build-cli:
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(CLI_BINARY_NAME) ./cmd/nango-cli

build-server:
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(SERVER_BINARY_NAME) ./cmd/nango-server

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	$(GOCLEAN)
	rm -rf $(BINARY_DIR)
	rm -f coverage.out coverage.html

install: build
	cp $(BINARY_DIR)/$(CLI_BINARY_NAME) $(GOPATH)/bin/
	cp $(BINARY_DIR)/$(SERVER_BINARY_NAME) $(GOPATH)/bin/

fmt:
	$(GOFMT) -s -w .

vet:
	$(GOVET) ./...

lint: fmt vet

deps:
	$(GOMOD) download
	$(GOMOD) tidy

run-cli: build-cli
	./$(BINARY_DIR)/$(CLI_BINARY_NAME) -help

run-server: build-server
	./$(BINARY_DIR)/$(SERVER_BINARY_NAME)

example: build
	$(GOBUILD) -o $(BINARY_DIR)/example ./examples/basic
	./$(BINARY_DIR)/example

# Docker targets (optional)
docker-build:
	docker build -t go-nango .

docker-run: docker-build
	docker run -p 8080:8080 go-nango

# Development helpers
dev-cli: build-cli
	@echo "Running CLI with environment variables..."
	@echo "Set NANGO_API_KEY before running:"
	@echo "export NANGO_API_KEY=your-api-key"
	@echo "Then run: ./$(BINARY_DIR)/$(CLI_BINARY_NAME) -command list"

dev-server: build-server
	@echo "Running server with environment variables..."
	@echo "Set required environment variables before running:"
	@echo "export NANGO_API_KEY=your-api-key"
	@echo "export PORT=8080"
	@echo "Then run: ./$(BINARY_DIR)/$(SERVER_BINARY_NAME)"