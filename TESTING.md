# Testing Guide for 4me Todos

This document covers all testing strategies, setup instructions, and examples for the 4me Todos application.

## Table of Contents

1. [Testing Overview](#testing-overview)
2. [Backend Testing (Go)](#backend-testing-go)
3. [Frontend Testing (Vue.js)](#frontend-testing-vuejs)
4. [E2E Testing](#e2e-testing)
5. [Running Tests](#running-tests)
6. [Test Coverage](#test-coverage)
7. [CI/CD Integration](#cicd-integration)

---

## Testing Overview

The application uses a comprehensive testing strategy:

- **Unit Tests**: Test individual functions and components in isolation
- **Integration Tests**: Test multiple components working together
- **E2E Tests**: Test complete user workflows from UI to database

### Testing Tools

**Backend (Go)**:

- `testing` - Go standard library
- `testify/assert` - Assertions and test helpers
- `testify/suite` - Test suite organization
- `httptest` - HTTP handler testing

**Frontend (Vue.js)**:

- `vitest` - Fast unit test framework
- `@vue/test-utils` - Vue component testing utilities
- `happy-dom` - Lightweight DOM implementation

---

## Backend Testing (Go)

### Test Structure

```
backend/
├── internal/
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── auth_test.go          # Unit tests for auth handler
│   │   ├── projects.go
│   │   └── projects_test.go      # Unit tests for projects
│   └── utils/
│       ├── jwt.go
│       └── jwt_test.go           # Unit tests for JWT utils
└── tests/
    └── e2e/
        └── api_test.go           # E2E integration tests
```

### Setting Up Test Database

Before running backend tests, create a test database:

```bash
# Create test database
createdb 4me_todos_test

# Or using psql
psql -U postgres -c "CREATE DATABASE 4me_todos_test;"
```

### Writing Unit Tests

Example: Testing the Auth Handler

```go
// internal/handlers/auth_test.go
package handlers

import (
 "bytes"
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "testing"

 "github.com/gin-gonic/gin"
 "github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
 // Setup
 db := setupTestDB(t)
 defer db.Close()

 handler := NewAuthHandler(db, testConfig)
 router := setupTestRouter(handler)

 // Test case
 t.Run("successful registration", func(t *testing.T) {
  reqBody := models.RegisterRequest{
   Username: "testuser",
   Email:    "test@example.com",
   Password: "password123",
  }

  body, _ := json.Marshal(reqBody)
  req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
  req.Header.Set("Content-Type", "application/json")
  w := httptest.NewRecorder()

  router.ServeHTTP(w, req)

  // Assertions
  assert.Equal(t, http.StatusCreated, w.Code)

  var response models.AuthResponse
  err := json.Unmarshal(w.Body.Bytes(), &response)
  assert.NoError(t, err)
  assert.NotEmpty(t, response.Token)
  assert.Equal(t, "testuser", response.User.Username)
 })
}
```

### Running Backend Tests

```bash
# Navigate to backend directory
cd backend

# Run all tests
make test

# Run unit tests only
make test-unit

# Run e2e tests only
make test-e2e

# Run with coverage
make test-coverage

# Run specific package
go test -v ./internal/handlers/...

# Run specific test
go test -v -run TestRegister ./internal/handlers/
```

### Test Coverage Report

```bash
cd backend
make test-coverage

# Opens coverage.html in browser
# Shows line-by-line coverage visualization
```

---

## Frontend Testing (Vue.js)

### Test Structure

```
frontend/
└── src/
    ├── components/
    │   └── __tests__/
    │       └── TaskCard.spec.js
    ├── stores/
    │   └── __tests__/
    │       ├── auth.spec.js
    │       └── projects.spec.js
    └── views/
        └── __tests__/
            └── DashboardView.spec.js
```

### Writing Component Tests

Example: Testing a Vue Component

```javascript
// src/components/__tests__/TaskCard.spec.js
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import TaskCard from '../TaskCard.vue'

describe('TaskCard', () => {
  it('renders task title', () => {
    const task = {
      id: 1,
      title: 'Test Task',
      priority: 'high',
      labels: []
    }

    const wrapper = mount(TaskCard, {
      props: { task }
    })

    expect(wrapper.text()).toContain('Test Task')
  })

  it('emits click event when clicked', async () => {
    const task = {
      id: 1,
      title: 'Test Task',
      priority: 'medium',
      labels: []
    }

    const wrapper = mount(TaskCard, {
      props: { task }
    })

    await wrapper.trigger('click')

    expect(wrapper.emitted('click')).toBeTruthy()
  })

  it('displays priority badge with correct color', () => {
    const task = {
      id: 1,
      title: 'Test Task',
      priority: 'urgent',
      labels: []
    }

    const wrapper = mount(TaskCard, {
      props: { task }
    })

    const badge = wrapper.find('.v-chip')
    expect(badge.exists()).toBe(true)
  })
})
```

### Testing Pinia Stores

Example: Testing Authentication Store

```javascript
// src/stores/__tests__/auth.spec.js
import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAuthStore } from '../auth'
import apiClient from '@/api/client'

// Mock the API client
vi.mock('@/api/client')

describe('Auth Store', () => {
  beforeEach(() => {
    setActivePinia(createPinia())
    localStorage.clear()
    vi.clearAllMocks()
  })

  it('successfully logs in user', async () => {
    const mockResponse = {
      data: {
        token: 'test-token',
        user: { id: 1, username: 'testuser' }
      }
    }
    
    apiClient.post.mockResolvedValue(mockResponse)
    
    const store = useAuthStore()
    const result = await store.login({
      username: 'testuser',
      password: 'password123'
    })
    
    expect(result.success).toBe(true)
    expect(store.token).toBe('test-token')
    expect(store.user.username).toBe('testuser')
  })
})
```

### Running Frontend Tests

```bash
# Navigate to frontend directory
cd frontend

# Install test dependencies
npm install -D vitest @vue/test-utils happy-dom

# Run tests
npm test

# Run tests in watch mode
npm test -- --watch

# Run tests with UI
npm run test:ui

# Run with coverage
npm run test:coverage
```

---

## E2E Testing

### Complete Workflow Tests

E2E tests verify entire user workflows from authentication to task management.

```go
// backend/tests/e2e/api_test.go
func (suite *E2ETestSuite) TestCompleteWorkflow() {
 // 1. Register a user
 registerReq := models.RegisterRequest{
  Username: "testuser",
  Email:    "test@example.com",
  Password: "password123",
 }
 // ... registration logic ...

 // 2. Create a project
 projectReq := models.CreateProjectRequest{
  Name: "My Todo Project",
 }
 // ... project creation ...

 // 3. Create boards
 boardReq := models.CreateBoardRequest{
  Name: "To Do",
 }
 // ... board creation ...

 // 4. Create tasks
 taskReq := models.CreateTaskRequest{
  Title:    "Fix login bug",
  Priority: "high",
 }
 // ... task creation ...

 // 5. Add comments
 // 6. Upload attachments
 // 7. Verify task history
}
```

### Test Data Management

```go
func (suite *E2ETestSuite) SetupTest() {
 // Clean database before each test
 suite.db.Pool.Exec(context.Background(), 
  "TRUNCATE users, projects, boards, tasks, labels, task_labels, comments, attachments, task_history CASCADE")
}
```

---

## Running Tests

### Complete Test Suite

```bash
# Backend (from project root)
cd backend
make test              # Run all backend tests
make test-coverage     # With coverage report

# Frontend (from project root)
cd frontend
npm test              # Run all frontend tests
npm run test:coverage # With coverage report
```

### Pre-Deployment Testing

```bash
# 1. Setup test database
createdb 4me_todos_test

# 2. Run backend tests
cd backend
make test

# 3. Build backend
make build

# 4. Run frontend tests
cd ../frontend
npm test

# 5. Build frontend
npm run build

# 6. Run e2e tests
cd ../backend
make test-e2e
```

### Continuous Testing During Development

```bash
# Terminal 1: Backend tests in watch mode
cd backend
# Go doesn't have native watch, use a tool like:
# go install github.com/cespare/reflex@latest
# reflex -r '\.go$' make test-unit

# Terminal 2: Frontend tests in watch mode
cd frontend
npm test -- --watch
```

---

## Test Coverage

### Coverage Goals

- **Backend**: Aim for >80% coverage on handlers and business logic
- **Frontend**: Aim for >70% coverage on stores and critical components

### Viewing Coverage Reports

**Backend**:

```bash
cd backend
make test-coverage
# Opens coverage.html in browser
```

**Frontend**:

```bash
cd frontend
npm run test:coverage
# Opens coverage/index.html in browser
```

### Coverage Reports Show

- **Line Coverage**: Percentage of code lines executed
- **Branch Coverage**: Percentage of conditional branches tested
- **Function Coverage**: Percentage of functions called
- **Statement Coverage**: Percentage of statements executed

---

## CI/CD Integration

### GitHub Actions Example

Create `.github/workflows/test.yml`:

```yaml
name: Tests

on: [push, pull_request]

jobs:
  backend-tests:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:14
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: 4me_todos_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Install dependencies
        working-directory: ./backend
        run: go mod download
      
      - name: Run tests
        working-directory: ./backend
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/4me_todos_test?sslmode=disable
        run: go test -v -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./backend/coverage.out

  frontend-tests:
    runs-on: ubuntu-latest
    
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node
        uses: actions/setup-node@v3
        with:
          node-version: '18'
      
      - name: Install dependencies
        working-directory: ./frontend
        run: npm ci
      
      - name: Run tests
        working-directory: ./frontend
        run: npm test -- --coverage
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          files: ./frontend/coverage/coverage-final.json
```

---

## Testing Best Practices

### 1. Test Naming Convention

```go
// Backend (Go)
func TestHandlerName_Scenario_ExpectedBehavior(t *testing.T) {}

// Example:
func TestAuthHandler_Register_SuccessfullyCreatesUser(t *testing.T) {}
func TestAuthHandler_Register_RejectsInvalidEmail(t *testing.T) {}
```

```javascript
// Frontend (JavaScript)
describe('ComponentName', () => {
  it('does something when condition', () => {})
})

// Example:
describe('TaskCard', () => {
  it('emits click event when clicked', () => {})
  it('displays priority badge with correct color', () => {})
})
```

### 2. AAA Pattern (Arrange, Act, Assert)

```go
func TestExample(t *testing.T) {
 // Arrange - Setup test data
 db := setupTestDB(t)
 handler := NewHandler(db)
 
 // Act - Execute the test
 result := handler.DoSomething()
 
 // Assert - Verify expectations
 assert.Equal(t, expected, result)
}
```

### 3. Test Isolation

- Each test should be independent
- Clean database between tests
- Don't rely on test execution order
- Use fresh test data for each test

### 4. Mock External Dependencies

```javascript
// Mock API calls
vi.mock('@/api/client', () => ({
  default: {
    post: vi.fn(),
    get: vi.fn()
  }
}))

// Mock stores
vi.mock('@/stores/auth', () => ({
  useAuthStore: vi.fn(() => ({
    user: { id: 1, username: 'test' },
    token: 'test-token'
  }))
}))
```

### 5. Test Edge Cases

- Empty inputs
- Null/undefined values
- Large datasets
- Concurrent requests
- Network failures
- Invalid tokens

---

## Troubleshooting Tests

### Common Issues

**Database Connection Errors**:

```bash
# Ensure PostgreSQL is running
pg_isadmin

# Create test database if missing
createdb 4me_todos_test

# Check DATABASE_URL in test files
```

**Frontend Test Failures**:

```bash
# Clear node_modules and reinstall
rm -rf node_modules package-lock.json
npm install

# Clear test cache
npm test -- --clearCache
```

**Import Errors**:

```bash
# Ensure path aliases are configured in vitest.config.js
resolve: {
  alias: {
    '@': fileURLToPath(new URL('./src', import.meta.url))
  }
}
```

---

## Next Steps

1. **Expand Test Coverage**: Add tests for remaining components
2. **Add Snapshot Testing**: For UI components
3. **Performance Testing**: Load test the API endpoints
4. **Security Testing**: Penetration testing, OWASP checks
5. **Accessibility Testing**: A11y compliance verification

For more information, see:

- [Architecture Documentation](./ARCHITECTURE.md)
- [Deployment Guide](./DEPLOYMENT.md)
- [Main README](./README.md)
