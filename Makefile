# BakaSub Makefile
# Cross-platform build system

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "v1.0.0-dev")
LDFLAGS := -s -w -X github.com/lsilvatti/bakasub/pkg/utils.Version=$(VERSION)
BINARY_NAME := bakasub
BUILD_DIR := bin

.PHONY: all clean build-linux build-windows build-macos build-all test install help

all: build-linux

help: ## Show this help message
	@echo "BakaSub - Build System"
	@echo "======================"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

clean: ## Remove all build artifacts
	@echo "🧹 Cleaning build directory..."
	@rm -rf $(BUILD_DIR)
	@echo "✓ Clean complete"

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test -v ./...
	@echo "✓ Tests complete"

build-linux: ## Build for Linux (AMD64)
	@echo "🐧 Building for Linux..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 \
		./cmd/$(BINARY_NAME)
	@echo "✓ Linux build complete: $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64"

build-windows: ## Build for Windows (AMD64)
	@echo "🪟 Building for Windows..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe \
		./cmd/$(BINARY_NAME)
	@echo "✓ Windows build complete: $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe"

build-macos: ## Build for macOS (AMD64 + ARM64)
	@echo "🍎 Building for macOS (Intel)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 \
		./cmd/$(BINARY_NAME)
	@echo "✓ macOS Intel build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64"
	@echo "🍎 Building for macOS (Apple Silicon)..."
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build \
		-ldflags="$(LDFLAGS)" \
		-o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 \
		./cmd/$(BINARY_NAME)
	@echo "✓ macOS ARM build complete: $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64"

build-all: clean build-linux build-windows build-macos ## Build for all platforms
	@echo ""
	@echo "🎉 All builds complete!"
	@echo ""
	@ls -lh $(BUILD_DIR)

install: build-linux ## Build and install to /usr/local/bin
	@echo "📦 Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 /usr/local/bin/$(BINARY_NAME)
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "✓ Installation complete!"
	@echo "Run '$(BINARY_NAME)' to start"

# Development targets
dev: ## Run in development mode
	@go run ./cmd/$(BINARY_NAME)

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✓ Format complete"

lint: ## Run linter
	@echo "🔍 Running linter..."
	@golangci-lint run ./... || echo "⚠ Install golangci-lint: https://golangci-lint.run/usage/install/"

# Release helper
release: ## Create a new release (use: make release VERSION=v1.2.3)
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ Error: VERSION not specified"; \
		echo "Usage: make release VERSION=v1.2.3"; \
		exit 1; \
	fi
	@echo "🚀 Creating release $(VERSION)..."
	@git tag -a $(VERSION) -m "Release $(VERSION)"
	@git push origin $(VERSION)
	@echo "✓ Release tag created and pushed"
	@echo "GitHub Actions will build and publish the release automatically"
