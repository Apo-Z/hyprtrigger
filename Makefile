.PHONY: build install clean test lint version help daemon reload status shutdown

APP_NAME = hyprtrigger
VERSION  ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "dev")
COMMIT    = $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE      = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

LDFLAGS = -X hyprtrigger/cmd.version=$(VERSION) \
          -X hyprtrigger/cmd.commit=$(COMMIT) \
          -X hyprtrigger/cmd.date=$(DATE)

all: build

build:
	@echo "Building $(APP_NAME) $(VERSION)..."
	go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)

install:
	go install -ldflags="$(LDFLAGS)"

clean:
	go clean
	rm -f $(APP_NAME)

test:
	go test ./...

lint:
	golangci-lint run

version:
	@echo "$(APP_NAME) $(VERSION) ($(COMMIT)) built on $(DATE)"

daemon: build
	./$(APP_NAME)

setup: build
	./$(APP_NAME) init-config

build-all:
	GOOS=linux  GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-linux-amd64
	GOOS=linux  GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(APP_NAME)-darwin-arm64

# Daemon control
reload:
	./$(APP_NAME) reload

status:
	./$(APP_NAME) status

shutdown:
	./$(APP_NAME) shutdown

# Development targets
dev-daemon:
	go run -ldflags="$(LDFLAGS)" .

dev-reload:
	go run -ldflags="$(LDFLAGS)" . reload

dev-status:
	go run -ldflags="$(LDFLAGS)" . status

dev-shutdown:
	go run -ldflags="$(LDFLAGS)" . shutdown

dev-config:
	go run -ldflags="$(LDFLAGS)" . init-config

dev-test:
	go run -ldflags="$(LDFLAGS)" . events list

dev-export:
	go run -ldflags="$(LDFLAGS)" . events export example.json

# Create new builtin event file
new-app:
	@if [ -z "$(APP)" ]; then \
		echo "Usage: make new-app APP=<app-name>"; \
		exit 1; \
	fi
	@chmod +x scripts/new-app.sh
	@./scripts/new-app.sh "$(APP)"

help:
	@./$(APP_NAME) --help
