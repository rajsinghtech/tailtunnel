.PHONY: help install build-frontend build-backend build dev run clean

# Load .env file if it exists
ifneq (,$(wildcard .env))
    include .env
    export
endif

# Default values (can be overridden by .env or environment)
TS_AUTHKEY ?=
STATE_DIR ?= $(PWD)/tailtunnel-state
PORT ?= 8080

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	cd frontend && npm install --legacy-peer-deps

build-frontend: ## Build the frontend
	@echo "Building frontend..."
	cd frontend && npm run build

build-backend: build-frontend ## Build the backend (includes frontend build)
	@echo "Building backend..."
	go build -o tailtunnel ./cmd/tailtunnel

build: build-backend ## Build everything

dev-frontend: ## Run frontend in development mode
	@echo "Starting frontend dev server..."
	cd frontend && DEVELOPMENT_BACKEND_URL=http://localhost:$(PORT) npm run dev

dev-backend: ## Run backend in development mode
	@echo "Starting backend dev server..."
	@mkdir -p $(STATE_DIR)
	TS_AUTHKEY=$(TS_AUTHKEY) STATE_DIR=$(STATE_DIR) PORT=$(PORT) go run ./cmd/tailtunnel

run: build ## Build and run the application
	@echo "Starting TailTunnel..."
	@mkdir -p $(STATE_DIR)
	TS_AUTHKEY=$(TS_AUTHKEY) STATE_DIR=$(STATE_DIR) PORT=$(PORT) ./tailtunnel

dev: ## Run in development mode (backend only, use 'make dev-frontend' in another terminal)
	@echo "Starting TailTunnel in development mode..."
	@echo "Backend will run on http://localhost:$(PORT)"
	@echo "Run 'make dev-frontend' in another terminal for frontend dev server"
	@echo ""
	@echo "Debug info:"
	@echo "  STATE_DIR=$(STATE_DIR)"
	@echo "  PORT=$(PORT)"
	@echo "  TS_AUTHKEY=$(TS_AUTHKEY)"
	@echo ""
	@mkdir -p $(STATE_DIR)
	TS_AUTHKEY=$(TS_AUTHKEY) STATE_DIR=$(STATE_DIR) PORT=$(PORT) go run ./cmd/tailtunnel

debug: ## Run with verbose logging
	@echo "Starting TailTunnel in debug mode..."
	@mkdir -p $(STATE_DIR)
	TS_AUTHKEY=$(TS_AUTHKEY) STATE_DIR=$(STATE_DIR) PORT=$(PORT) TS_DEBUG=true go run -v ./cmd/tailtunnel

logs: ## Show Tailscale logs from state directory
	@if [ -d "$(STATE_DIR)" ]; then \
		echo "Tailscale state directory contents:"; \
		ls -lah $(STATE_DIR); \
	else \
		echo "State directory not found: $(STATE_DIR)"; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -rf tailtunnel
	rm -rf frontend/dist
	rm -rf frontend/.svelte-kit
	rm -rf frontend/node_modules
	rm -rf $(STATE_DIR)

test: ## Run tests
	@echo "Running Go tests..."
	go test ./...

tidy: ## Tidy Go modules
	@echo "Tidying Go modules..."
	go mod tidy

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t tailtunnel:latest .

docker-run: ## Run Docker container
	@echo "Running Docker container..."
	@if [ -z "$(TS_AUTHKEY)" ]; then \
		echo "Error: TS_AUTHKEY not set. Please create a .env file or set the environment variable."; \
		exit 1; \
	fi
	docker run --rm -it \
		--cap-add=NET_ADMIN \
		-p $(PORT):8080 \
		-e TS_AUTHKEY=$(TS_AUTHKEY) \
		-e PORT=8080 \
		-v tailtunnel-state:/var/lib/tailtunnel \
		tailtunnel:latest

docker-compose-up: ## Start with docker-compose
	@echo "Starting with docker-compose..."
	@if [ -z "$(TS_AUTHKEY)" ]; then \
		echo "Error: TS_AUTHKEY not set. Please create a .env file."; \
		exit 1; \
	fi
	docker-compose up -d

docker-compose-down: ## Stop docker-compose services
	@echo "Stopping docker-compose services..."
	docker-compose down

docker-compose-logs: ## Show docker-compose logs
	docker-compose logs -f

docker-clean: ## Clean Docker images and volumes
	@echo "Cleaning Docker artifacts..."
	docker-compose down -v
	docker rmi tailtunnel:latest 2>/dev/null || true

build-macos-app: build-frontend ## Build the macOS menu bar app
	@echo "Building macOS app..."
	@mkdir -p TailTunnel.app/Contents/{MacOS,Resources}
	@cp resources/Info.plist TailTunnel.app/Contents/
	@cp resources/TailTunnel.icns TailTunnel.app/Contents/Resources/
	@go build -o TailTunnel.app/Contents/MacOS/TailTunnel -ldflags="-s -w" ./cmd/tailtunnel-menubar
	@echo "macOS app built at: TailTunnel.app"
	@echo "To install: cp -r TailTunnel.app /Applications/"

.DEFAULT_GOAL := help
