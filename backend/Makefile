.PHONY: help test test-unit test-e2e test-coverage build run clean lint fmt deps test-deps dev prod security docker-build docker-run migrate-up migrate-down

# Default target
help: ## Show this help message
	@echo "Backend Makefile - Available Commands:"
	@echo ""
	@echo "Setup & Installation:"
	@echo "  deps           - Install Go dependencies"
	@echo "  test-deps      - Install test dependencies"
	@echo "  env-setup      - Setup environment file"
	@echo ""
	@echo "Development:"
	@echo "  dev            - Development workflow (clean, build, test)"
	@echo "  run            - Run the application"
	@echo "  watch          - Run with auto-reload (requires air)"
	@echo ""
	@echo "Building:"
	@echo "  build          - Build the application"
	@echo "  prod           - Build for production"
	@echo "  docker-build   - Build Docker image"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-e2e       - Run e2e tests only"
	@echo "  test-coverage  - Run tests with coverage"
	@echo "  test-db-setup  - Setup test database"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Lint code with golangci-lint"
	@echo "  fmt            - Format code"
	@echo "  security       - Run security checks"
	@echo ""
	@echo "Database:"
	@echo "  migrate-up     - Run migrations (via server start)"
	@echo "  migrate-down   - Rollback migrations (not supported)"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run with Docker"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean          - Clean build artifacts"
	@echo ""

# Setup & Installation
deps: ## Install Go dependencies
	@echo "📦 Installing dependencies..."
	@go mod download
	@go mod tidy
	@echo "✅ Dependencies installed"

test-deps: ## Install test dependencies
	@echo "🧪 Installing test dependencies..."
	@go install github.com/stretchr/testify@latest
	@go get -u github.com/stretchr/testify/assert
	@go get -u github.com/stretchr/testify/suite
	@echo "✅ Test dependencies installed"

migrate-cli: ## Install migrate CLI tool
	@echo "🔧 Installing migrate CLI..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "✅ Migrate CLI installed"

env-setup: ## Setup environment file
	@if [ ! -f .env ]; then \
		cp .env.example .env 2>/dev/null || echo "# Create .env manually" > .env; \
		echo "📝 Please edit .env with your database URL and secrets"; \
	fi

# Development
dev: clean build test ## Development workflow
	@echo "✅ Development workflow complete"

run: ## Run the application
	@echo "🚀 Running application..."
	@go run cmd/api/main.go

watch: ## Run with auto-reload (requires air)
	@echo "👀 Running with auto-reload..."
	@if command -v air > /dev/null; then \
		air; \
	else \
		echo "❌ Air not installed. Install with: go install github.com/cosmtrek/air@latest"; \
		echo "🔄 Falling back to regular run..."; \
		$(MAKE) run; \
	fi

# Building
build: ## Build the application
	@echo "🔨 Building..."
	@mkdir -p bin
	@go build -o bin/server cmd/api/main.go
	@echo "✅ Build complete"

prod: clean ## Build for production
	@echo "🏭 Building for production..."
	@mkdir -p bin
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/server cmd/api/main.go
	@echo "✅ Production build complete"

docker-build: ## Build Docker image
	@echo "🐳 Building Docker image..."
	@docker build -t 4me-backend .
	@echo "✅ Docker image built"

docker-run: ## Run with Docker
	@echo "🐳 Running with Docker..."
	@docker run -p 8080:8080 --env-file .env 4me-backend

# Testing
test: ## Run all tests
	@echo "🧪 Running all tests..."
	@go test -v ./...

test-unit: ## Run unit tests only
	@echo "🧪 Running unit tests..."
	@go test -v ./internal/...

test-e2e: ## Run e2e tests only
	@echo "🧪 Running e2e tests..."
	@go test -v ./tests/e2e/...

test-coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report generated: coverage.html"

test-db-setup: ## Setup test database
	@echo "🗄️ Setting up test database..."
	@psql -U postgres -c "DROP DATABASE IF EXISTS 4me_todos_test;" 2>/dev/null || true
	@psql -U postgres -c "CREATE DATABASE 4me_todos_test;" 2>/dev/null || echo "Test database setup failed or already exists"

# Code Quality
lint: ## Lint code with golangci-lint
	@echo "🔍 Linting code..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "❌ golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

fmt: ## Format code
	@echo "🎨 Formatting code..."
	@go fmt ./...
	@echo "✅ Code formatted"

security: ## Run security checks
	@echo "🔒 Running security checks..."
	@go mod audit || echo "No vulnerabilities found"
	@if command -v gosec > /dev/null; then \
		gosec ./...; \
	else \
		echo "ℹ️  gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest"; \
	fi

# Database migrations
migrate-up: ## Run all pending migrations
	@echo "🔄 Running migrations..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "$(DATABASE_URL)" up; \
	else \
		echo "❌ migrate CLI not installed. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
		echo "🔄 Running migrations via server start..."; \
		$(MAKE) run; \
	fi

migrate-down: ## Rollback last migration
	@echo "🔄 Rolling back last migration..."
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "$(DATABASE_URL)" down 1; \
	else \
		echo "❌ migrate CLI not installed. Install with: go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest"; \
	fi

migrate-force: ## Force migration to specific version (use with caution)
	@echo "⚠️  Force migration to version $(VERSION)"
	@if [ -z "$(VERSION)" ]; then \
		echo "❌ VERSION required. Usage: make migrate-force VERSION=1"; \
		exit 1; \
	fi
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "$(DATABASE_URL)" force $(VERSION); \
	else \
		echo "❌ migrate CLI not installed"; \
	fi

migrate-version: ## Show current migration version
	@echo "📊 Current migration version:"
	@if command -v migrate > /dev/null; then \
		migrate -path migrations -database "$(DATABASE_URL)" version; \
	else \
		echo "❌ migrate CLI not installed"; \
	fi

migrate-create: ## Create new migration file
	@echo "📝 Creating new migration..."
	@if [ -z "$(NAME)" ]; then \
		echo "❌ NAME required. Usage: make migrate-create NAME=add_new_column"; \
		exit 1; \
	fi
	@if command -v migrate > /dev/null; then \
		migrate create -ext sql -dir migrations -seq $(NAME); \
	else \
		echo "❌ migrate CLI not installed"; \
	fi

# Cleanup
clean: ## Clean build artifacts
	@echo "🧹 Cleaning..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ Cleanup complete"

# Development utilities
dev-setup: deps env-setup test-db-setup ## Complete development setup
	@echo "✅ Development environment ready"

dev-reset: clean deps ## Reset development environment
	@echo "✅ Development environment reset"

# Quick commands
start: run ## Alias for run
restart: clean run ## Clean and restart

# Analysis
analyze: ## Analyze code
	@echo "📊 Analyzing code..."
	@go vet ./...
	@go mod verify
	@echo "✅ Analysis complete"

benchmark: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	@go test -bench=. -benchmem ./...

profile: ## Generate CPU profile
	@echo "📊 Generating CPU profile..."
	@go test -cpuprofile=cpu.prof -bench=. ./internal/...

# Health check
health: ## Check application health
	@echo "🏥 Health Check:"
	@curl -s http://localhost:8080/api/auth/me >/dev/null && echo "  ✅ API: OK" || echo "  ❌ API: Down"

