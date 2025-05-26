.PHONY: build install clean test lint version help run reload status shutdown

APP_NAME = hyprtrigger
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
COMMIT = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -X hyprtrigger/cmd/hyprtrigger.version=$(VERSION) -X hyprtrigger/cmd/hyprtrigger.commit=$(COMMIT) -X hyprtrigger/cmd/hyprtrigger.date=$(DATE)

all: build

build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)

install:
	@echo "Installing $(APP_NAME) $(VERSION)..."
	go install -ldflags="$(LDFLAGS)"

clean:
	@echo "Cleaning..."
	go clean
	rm -f $(APP_NAME)

test:
	@echo "Running tests..."
	go test ./...

lint:
	@echo "Running linter..."
	golangci-lint run

version:
	@echo "$(APP_NAME) $(VERSION) ($(COMMIT)) built on $(DATE)"

run: build
	./$(APP_NAME) --daemon

daemon: build
	./$(APP_NAME) --daemon

setup: build
	./$(APP_NAME) --create-config

build-all:
	@echo "Building for multiple platforms..."
	GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-darwin-arm64

# Development targets
dev-run:
	go run -ldflags="$(LDFLAGS)" . --daemon

dev-daemon:
	go run -ldflags="$(LDFLAGS)" . --daemon

dev-config:
	go run -ldflags="$(LDFLAGS)" . --create-config
dev-test:
	go run -ldflags="$(LDFLAGS)" . --print-builtin

dev-export:
	go run -ldflags="$(LDFLAGS)" . --export-builtin example.json

# Daemon control targets
reload:
	@echo "Reloading hyprtrigger configuration..."
	./$(APP_NAME) --reload

status:
	@echo "Checking hyprtrigger daemon status..."
	./$(APP_NAME) --status

shutdown:
	@echo "Shutting down hyprtrigger daemon..."
	./$(APP_NAME) --shutdown

# Development reload workflow
dev-reload:
	go run -ldflags="$(LDFLAGS)" . --reload

dev-status:
	go run -ldflags="$(LDFLAGS)" . --status

dev-shutdown:
	go run -ldflags="$(LDFLAGS)" . --shutdown

# Create new builtin event file
new-app:
	@if [ -z "$(APP)" ]; then \
		echo "Usage: make new-app APP=<app-name>"; \
		echo "Example: make new-app APP=firefox"; \
		exit 1; \
	fi
	@chmod +x scripts/new-app.sh
	@./scripts/new-app.sh "$(APP)"

help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  install      - Install the application"
	@echo "  clean        - Clean build artifacts"
	@echo "  test         - Run tests"
	@echo "  lint         - Run linter"
	@echo "  version      - Show version information"
	@echo "  run          - Build and start daemon"
	@echo "  daemon       - Build and start daemon (alias for run)"
	@echo "  setup        - Build and create config directory"
	@echo "  build-all    - Build for multiple platforms"
	@echo ""
	@echo "Development targets:"
	@echo "  dev-run      - Start daemon (development)"
	@echo "  dev-daemon   - Start daemon (development, alias)"
	@echo "  dev-config   - Create config directory (development)"
	@echo "  dev-test     - Print builtin events (development)"
	@echo "  dev-export   - Export builtin events to example.json"
	@echo "  dev-reload   - Send reload command (development)"
	@echo "  dev-status   - Check daemon status (development)"
	@echo "  dev-shutdown - Shutdown daemon (development)"
	@echo "  new-app      - Create new builtin app (make new-app APP=name)"
	@echo ""
	@echo "Daemon control:"
	@echo "  reload       - Reload configuration in running daemon"
	@echo "  status       - Check daemon status"
	@echo "  shutdown     - Shutdown running daemon"
	@echo ""
	@echo "Hot reload workflow:"
	@echo "  1. make daemon                 # Start daemon"
	@echo "  2. Edit ~/.config/hyprtrigger/*.json"
	@echo "  3. make reload                 # Hot reload config"
	@echo "  4. make status                 # Check status"
	@echo "  5. make shutdown               # Stop daemon"
	@echo ""
	@echo "  help         - Show this help"
