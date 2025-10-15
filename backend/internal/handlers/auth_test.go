package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/config"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *database.Database {
	cfg := &config.Config{
		DatabaseURL: "postgres://localhost:5432/4me_todos_test?sslmode=disable",
	}

	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		t.Skipf("Skipping test: database not available - %v", err)
	}

	// Clean tables
	_, _ = db.Pool.Exec(context.Background(), "TRUNCATE users, projects, boards, tasks, labels, task_labels, comments, attachments, task_history CASCADE")

	return db
}

func setupTestRouter(handler *AuthHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)

	return router
}

func TestRegister(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	handler := NewAuthHandler(db, cfg)
	router := setupTestRouter(handler)

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

		assert.Equal(t, http.StatusCreated, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
		assert.Equal(t, "testuser", response.User.Username)
	})

	t.Run("duplicate username", func(t *testing.T) {
		reqBody := models.RegisterRequest{
			Username: "testuser",
			Email:    "another@example.com",
			Password: "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
	})

	t.Run("invalid email", func(t *testing.T) {
		reqBody := map[string]string{
			"username": "newuser",
			"email":    "invalid-email",
			"password": "password123",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestLogin(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	cfg := &config.Config{
		JWTSecret: "test-secret",
	}

	handler := NewAuthHandler(db, cfg)
	router := setupTestRouter(handler)

	// Register a user first
	registerReq := models.RegisterRequest{
		Username: "loginuser",
		Email:    "login@example.com",
		Password: "password123",
	}
	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Run("successful login", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Username: "loginuser",
			Password: "password123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.AuthResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response.Token)
	})

	t.Run("invalid password", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Username: "loginuser",
			Password: "wrongpassword",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("non-existent user", func(t *testing.T) {
		loginReq := models.LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		}

		body, _ := json.Marshal(loginReq)
		req := httptest.NewRequest("POST", "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})
}
