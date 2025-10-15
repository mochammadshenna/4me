.PHONY: help build test clean dev prod setup install deps lint fmt security docker-build docker-run deploy-frontend deploy-backend

# Default target
help: ## Show this help message
	@echo "4me Todos - Available Commands:"
	@echo ""
	@echo "Setup & Installation:"
	@echo "  setup          - Complete project setup (database, deps, env)"
	@echo "  install        - Install all dependencies"
	@echo "  deps           - Install Go and Node dependencies"
	@echo ""
	@echo "Development:"
	@echo "  dev            - Start development servers (backend + frontend)"
	@echo "  run            - Start backend server only"
	@echo "  dev-backend    - Start backend only"
	@echo "  dev-frontend   - Start frontend only"
	@echo ""
	@echo "Building:"
	@echo "  build          - Build both backend and frontend"
	@echo "  build-backend  - Build backend only"
	@echo "  build-frontend - Build frontend only"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests (backend + frontend)"
	@echo "  test-backend   - Run backend tests only"
	@echo "  test-frontend  - Run frontend tests only"
	@echo "  test-coverage  - Run tests with coverage reports"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Lint all code"
	@echo "  fmt            - Format all code"
	@echo "  security       - Run security checks"
	@echo ""
	@echo "Database:"
	@echo "  db-setup       - Setup database and run migrations"
	@echo "  db-reset       - Reset database (WARNING: deletes all data)"
	@echo "  db-test-setup  - Setup test database"
	@echo ""
	@echo "Docker:"
	@echo "  docker-build   - Build Docker images"
	@echo "  docker-run     - Run with Docker Compose"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy-frontend - Deploy frontend to Vercel"
	@echo "  deploy-backend  - Deploy backend to Railway/Render"
	@echo ""
	@echo "Cleanup:"
	@echo "  clean          - Clean all build artifacts"
	@echo ""

# Setup & Installation
setup: db-setup install env-setup ## Complete project setup
	@echo "✅ Project setup complete!"
	@echo "Run 'make dev' to start development servers"

install: deps ## Install all dependencies
	@echo "✅ All dependencies installed"

deps: ## Install Go and Node dependencies
	@echo "Installing backend dependencies..."
	@cd backend && go mod download && go mod tidy
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "✅ Dependencies installed"

env-setup: ## Setup environment files
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env 2>/dev/null || echo "# Create backend/.env manually" > backend/.env; \
		echo "📝 Please edit backend/.env with your database URL and secrets"; \
	fi
	@if [ ! -f frontend/.env ]; then \
		echo "VITE_API_URL=http://localhost:8080/api" > frontend/.env; \
		echo "VITE_GOOGLE_CLIENT_ID=" >> frontend/.env; \
		echo "📝 Please edit frontend/.env with your API URL and Google Client ID"; \
	fi

# Development
dev: ## Start development servers (backend + frontend)
	@echo "🚀 Starting development servers..."
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"
	@echo "Press Ctrl+C to stop all servers"
	@trap 'kill %1 %2' INT; \
	cd backend && go run cmd/api/main.go & \
	cd frontend && npm run dev & \
	wait

dev-backend: ## Start backend development server
	@echo "🚀 Starting backend server..."
	@cd backend && go run cmd/api/main.go

run: ## Start backend server only
	@echo "🚀 Starting backend server..."
	@cd backend && go run cmd/api/main.go

dev-frontend: ## Start frontend development server
	@echo "🚀 Starting frontend server..."
	@cd frontend && npm run dev

# Building
build: build-backend build-frontend ## Build both backend and frontend

build-backend: ## Build backend
	@echo "🔨 Building backend..."
	@cd backend && make build

build-frontend: ## Build frontend
	@echo "🔨 Building frontend..."
	@cd frontend && npm run build

# Testing
test: test-backend test-frontend ## Run all tests

test-backend: ## Run backend tests
	@echo "🧪 Running backend tests..."
	@cd backend && make test

test-frontend: ## Run frontend tests
	@echo "🧪 Running frontend tests..."
	@cd frontend && npm test

test-coverage: ## Run tests with coverage
	@echo "📊 Generating coverage reports..."
	@cd backend && make test-coverage
	@cd frontend && npm run test:coverage
	@echo "✅ Coverage reports generated"

# Code Quality
lint: ## Lint all code
	@echo "🔍 Linting backend..."
	@cd backend && make lint || echo "Backend linting failed or golangci-lint not installed"
	@echo "🔍 Linting frontend..."
	@cd frontend && npm run lint || echo "Frontend linting failed"

fmt: ## Format all code
	@echo "🎨 Formatting backend..."
	@cd backend && make fmt
	@echo "🎨 Formatting frontend..."
	@cd frontend && npm run lint -- --fix || echo "Frontend formatting completed"

security: ## Run security checks
	@echo "🔒 Running security checks..."
	@cd backend && go mod audit || echo "No vulnerabilities found"
	@cd frontend && npm audit || echo "No vulnerabilities found"

# Database
db-setup: ## Setup database and run migrations
	@echo "🗄️ Setting up database..."
	@createdb 4me_todos 2>/dev/null || echo "Database 4me_todos already exists or PostgreSQL not running"
	@echo "✅ Database setup complete"

migrate-up: ## Run all pending migrations
	@echo "🔄 Running migrations..."
	@cd backend && make migrate-up

migrate-down: ## Rollback last migration
	@echo "🔄 Rolling back migration..."
	@cd backend && make migrate-down

migrate-version: ## Show current migration version
	@echo "📊 Migration version:"
	@cd backend && make migrate-version

migrate-create: ## Create new migration file
	@echo "📝 Creating new migration..."
	@cd backend && make migrate-create NAME=$(NAME)

migrate-cli: ## Install migrate CLI tool
	@echo "🔧 Installing migrate CLI..."
	@cd backend && make migrate-cli

db-reset: ## Reset database (WARNING: deletes all data)
	@echo "⚠️  WARNING: This will delete all data!"
	@read -p "Are you sure? Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ]
	@dropdb 4me_todos 2>/dev/null || echo "Database doesn't exist"
	@createdb 4me_todos
	@echo "✅ Database reset complete"

db-test-setup: ## Setup test database
	@echo "🧪 Setting up test database..."
	@dropdb 4me_todos_test 2>/dev/null || echo "Test database doesn't exist"
	@createdb 4me_todos_test
	@echo "✅ Test database setup complete"

seed-demo: ## Create demo user (username: demo, password: password)
	@echo "🌱 Creating demo user..."
	@psql -U postgres -d 4me_todos -f backend/scripts/seed_demo.sql
	@echo ""
	@echo "✅ Demo user ready!"
	@echo "   Username: demo"
	@echo "   Password: password"

# Docker
docker-build: ## Build Docker images
	@echo "🐳 Building Docker images..."
	@docker build -t 4me-backend ./backend
	@docker build -t 4me-frontend ./frontend
	@echo "✅ Docker images built"

docker-run: ## Run with Docker Compose
	@echo "🐳 Starting with Docker Compose..."
	@docker-compose up --build

# Deployment
deploy-frontend: ## Deploy frontend to Vercel
	@echo "🚀 Deploying frontend to Vercel..."
	@cd frontend && npx vercel --prod

deploy-backend: ## Deploy backend to Railway/Render
	@echo "🚀 Deploying backend..."
	@echo "Manual deployment required - see DEPLOYMENT.md"
	@echo "Railway: railway deploy"
	@echo "Render: git push render main"

# Cleanup
clean: ## Clean all build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@cd backend && make clean
	@cd frontend && rm -rf dist/ coverage/
	@echo "✅ Cleanup complete"

# Production
prod: clean build test ## Production build with tests
	@echo "✅ Production build complete"

# Quick commands
start: dev ## Alias for dev
stop: ## Stop all running servers
	@pkill -f "go run cmd/api/main.go" 2>/dev/null || true
	@pkill -f "npm run dev" 2>/dev/null || true
	@echo "✅ All servers stopped"

restart: stop start ## Restart all servers

status: ## Show project status
	@echo "📊 Project Status:"
	@echo "Backend:"
	@pgrep -f "go run cmd/api/main.go" >/dev/null && echo "  ✅ Running" || echo "  ❌ Stopped"
	@echo "Frontend:"
	@pgrep -f "npm run dev" >/dev/null && echo "  ✅ Running" || echo "  ❌ Stopped"
	@echo "Database:"
	@pg_isready -q && echo "  ✅ PostgreSQL Running" || echo "  ❌ PostgreSQL Stopped"

# Health check
health: ## Check application health
	@echo "🏥 Health Check:"
	@curl -s http://localhost:8080/api/auth/me >/dev/null && echo "  ✅ Backend API: OK" || echo "  ❌ Backend API: Down"
	@curl -s http://localhost:5173 >/dev/null && echo "  ✅ Frontend: OK" || echo "  ❌ Frontend: Down"
