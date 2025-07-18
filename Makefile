.PHONY: build test clean run install lint fmt vet help

# Variables
BINARY_NAME=claude-telegram-notifier
MAIN_PATH=./cmd
BUILD_DIR=./bin

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

build: ## Build the application
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

test: ## Run all tests
	@echo "Running tests..."
	@go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests with verbose output..."
	@go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	@go test -cover ./...

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean

run: build ## Build and run the application
	@echo "Running $(BINARY_NAME)..."
	@$(BUILD_DIR)/$(BINARY_NAME)

chat-client: build ## Run chat client to discover chat ID
	@echo "Starting Telegram chat ID discovery..."
	@$(BUILD_DIR)/$(BINARY_NAME) chat-client

install: ## Install dependencies
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

install-cli: build ## Install CLI binary to /usr/local/bin
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
	@sudo chmod +x /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) installed successfully!"
	@echo "You can now run: $(BINARY_NAME) <hook-type>"

uninstall-cli: ## Uninstall CLI binary from /usr/local/bin
	@echo "Uninstalling $(BINARY_NAME) from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)
	@echo "$(BINARY_NAME) uninstalled successfully!"

lint: ## Run golint
	@echo "Running lint..."
	@which golint > /dev/null || go install golang.org/x/lint/golint@latest
	@golint ./...

fmt: ## Format code
	@echo "Formatting code..."
	@go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	@go vet ./...

check: fmt vet lint test ## Run all checks (fmt, vet, lint, test)

# Development targets
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	@go mod download
	@go install golang.org/x/lint/golint@latest

# Chat ID discovery
discover-chat-id: chat-client ## Alias for chat-client

# Usage examples
example-user-prompt: build ## Run example with user-prompt-submit hook
	@echo '{"event_type":"user-prompt-submit","timestamp":"2024-01-01T12:00:00Z","message":"Test message","tool_used":"Read","file_path":"/path/to/file.go"}' | $(BUILD_DIR)/$(BINARY_NAME) user-prompt-submit

example-tool-call: build ## Run example with tool-call hook
	@echo '{"event_type":"tool-call","timestamp":"2024-01-01T12:00:00Z","tool_name":"Bash","file_path":"/path/to/script.sh","success":true}' | $(BUILD_DIR)/$(BINARY_NAME) tool-call

example-file-write: build ## Run example with file-write hook
	@echo '{"event_type":"file-write","timestamp":"2024-01-01T12:00:00Z","file_path":"/path/to/file.go","success":true,"lines_added":10,"lines_deleted":5}' | $(BUILD_DIR)/$(BINARY_NAME) file-write

example-session-start: build ## Run example with session-start hook
	@echo '{"event_type":"session-start","timestamp":"2024-01-01T12:00:00Z","working_directory":"/home/user/project","git_repo":true}' | $(BUILD_DIR)/$(BINARY_NAME) session-start

example-session-end: build ## Run example with session-end hook
	@echo '{"event_type":"session-end","timestamp":"2024-01-01T12:00:00Z","duration":"30m","files_modified":5}' | $(BUILD_DIR)/$(BINARY_NAME) session-end

example-notification: build ## Run example with notification hook
	@echo '{"event_type":"notification","timestamp":"2024-01-01T12:00:00Z","message":"Test notification message","level":"info","source":"claude-code"}' | $(BUILD_DIR)/$(BINARY_NAME) notification

# Docker targets (optional)
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	@docker build -t $(BINARY_NAME) .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@docker run --rm -e TELEGRAM_BOT_TOKEN -e TELEGRAM_CHAT_ID $(BINARY_NAME)