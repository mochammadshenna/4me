.PHONY: help install build dev test test-watch test-coverage lint fmt clean preview serve deploy

# Default target
help: ## Show this help message
	@echo "Frontend Makefile - Available Commands:"
	@echo ""
	@echo "Setup:"
	@echo "  install        - Install dependencies"
	@echo "  clean          - Clean build artifacts and node_modules"
	@echo ""
	@echo "Development:"
	@echo "  dev            - Start development server"
	@echo "  preview        - Preview production build"
	@echo "  serve          - Serve production build"
	@echo ""
	@echo "Building:"
	@echo "  build          - Build for production"
	@echo "  build-dev      - Build for development"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run tests"
	@echo "  test-watch     - Run tests in watch mode"
	@echo "  test-ui        - Run tests with UI"
	@echo "  test-coverage  - Run tests with coverage"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint           - Lint code"
	@echo "  fmt            - Format code (fix linting issues)"
	@echo ""
	@echo "Deployment:"
	@echo "  deploy         - Deploy to Vercel"
	@echo "  deploy-preview - Deploy preview to Vercel"
	@echo ""
	@echo "Analysis:"
	@echo "  analyze        - Analyze bundle size"
	@echo "  audit          - Audit dependencies"
	@echo ""

# Setup
install: ## Install dependencies
	@echo "📦 Installing dependencies..."
	@npm install
	@echo "✅ Dependencies installed"

clean: ## Clean build artifacts and node_modules
	@echo "🧹 Cleaning..."
	@rm -rf dist/
	@rm -rf coverage/
	@rm -rf node_modules/
	@rm -f package-lock.json
	@echo "✅ Cleanup complete"

# Development
dev: ## Start development server
	@echo "🚀 Starting development server..."
	@echo "URL: http://localhost:5173"
	@npm run dev

preview: build ## Preview production build
	@echo "👀 Previewing production build..."
	@npm run preview

serve: build ## Serve production build
	@echo "🌐 Serving production build..."
	@npx serve dist -s

# Building
build: ## Build for production
	@echo "🔨 Building for production..."
	@npm run build
	@echo "✅ Production build complete"

build-dev: ## Build for development
	@echo "🔨 Building for development..."
	@npm run build -- --mode development
	@echo "✅ Development build complete"

# Testing
test: ## Run tests
	@echo "🧪 Running tests..."
	@npm test

test-watch: ## Run tests in watch mode
	@echo "👀 Running tests in watch mode..."
	@npm test -- --watch

test-ui: ## Run tests with UI
	@echo "🎨 Running tests with UI..."
	@npm run test:ui

test-coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	@npm run test:coverage
	@echo "✅ Coverage report generated: coverage/index.html"

# Code Quality
lint: ## Lint code
	@echo "🔍 Linting code..."
	@npm run lint

fmt: ## Format code (fix linting issues)
	@echo "🎨 Formatting code..."
	@npm run lint -- --fix
	@echo "✅ Code formatted"

# Deployment
deploy: ## Deploy to Vercel
	@echo "🚀 Deploying to Vercel..."
	@npx vercel --prod
	@echo "✅ Deployed to production"

deploy-preview: ## Deploy preview to Vercel
	@echo "👀 Deploying preview to Vercel..."
	@npx vercel
	@echo "✅ Preview deployed"

# Analysis
analyze: ## Analyze bundle size
	@echo "📊 Analyzing bundle size..."
	@npm run build
	@npx vite-bundle-analyzer dist/assets/*.js
	@echo "✅ Bundle analysis complete"

audit: ## Audit dependencies
	@echo "🔒 Auditing dependencies..."
	@npm audit
	@echo "✅ Audit complete"

audit-fix: ## Fix audit issues
	@echo "🔧 Fixing audit issues..."
	@npm audit fix
	@echo "✅ Audit issues fixed"

# Development workflow
dev-setup: install ## Setup development environment
	@echo "⚙️ Setting up development environment..."
	@if [ ! -f .env ]; then \
		echo "VITE_API_URL=http://localhost:8080/api" > .env; \
		echo "VITE_GOOGLE_CLIENT_ID=" >> .env; \
		echo "📝 Created .env file - please edit with your settings"; \
	fi
	@echo "✅ Development environment ready"

dev-check: ## Check development environment
	@echo "🔍 Checking development environment..."
	@echo "Node version: $(shell node --version)"
	@echo "NPM version: $(shell npm --version)"
	@echo "Dependencies: $(shell npm list --depth=0 2>/dev/null | wc -l) packages installed"
	@if [ -f .env ]; then echo "✅ .env file exists"; else echo "❌ .env file missing"; fi
	@echo "✅ Environment check complete"

# Quick commands
start: dev ## Alias for dev
restart: clean install dev ## Clean install and start dev server

# Build variants
build-analyze: build analyze ## Build and analyze bundle

# Test variants
test-all: test test-coverage ## Run all tests with coverage
test-ci: ## Run tests for CI
	@echo "🔄 Running tests for CI..."
	@npm test -- --run --reporter=verbose

# Utility commands
info: ## Show project information
	@echo "📋 Project Information:"
	@echo "Name: $(shell node -p "require('./package.json').name")"
	@echo "Version: $(shell node -p "require('./package.json').version")"
	@echo "Node: $(shell node --version)"
	@echo "NPM: $(shell npm --version)"
	@echo "Dependencies: $(shell npm list --depth=0 2>/dev/null | grep -c '@')"

deps-update: ## Update dependencies
	@echo "🔄 Updating dependencies..."
	@npm update
	@echo "✅ Dependencies updated"

deps-outdated: ## Check for outdated dependencies
	@echo "📅 Checking for outdated dependencies..."
	@npm outdated
	@echo "✅ Outdated check complete"
