package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupProjectRouter(handler *ProjectHandler, userID int) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Mock auth middleware
	router.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})

	router.POST("/projects", handler.Create)
	router.GET("/projects", handler.List)
	router.GET("/projects/:id", handler.Get)
	router.PUT("/projects/:id", handler.Update)
	router.DELETE("/projects/:id", handler.Delete)

	return router
}

func createTestUser(t *testing.T, db *database.Database) int {
	var userID int
	err := db.Pool.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING id`,
		"testuser", "test@example.com", "hashedpassword").Scan(&userID)
	assert.NoError(t, err)
	return userID
}

func TestProjectCreate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	handler := NewProjectHandler(db)
	router := setupProjectRouter(handler, userID)

	t.Run("create project successfully", func(t *testing.T) {
		reqBody := models.CreateProjectRequest{
			Name:        "Test Project",
			Description: stringPtr("A test project"),
			Color:       "#3B82F6",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var project models.Project
		err := json.Unmarshal(w.Body.Bytes(), &project)
		assert.NoError(t, err)
		assert.Equal(t, "Test Project", project.Name)
		assert.Equal(t, "#3B82F6", project.Color)
	})

	t.Run("create project without name", func(t *testing.T) {
		reqBody := map[string]string{
			"description": "Test",
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("POST", "/projects", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestProjectList(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	handler := NewProjectHandler(db)
	router := setupProjectRouter(handler, userID)

	// Create test projects
	_, _ = db.Pool.Exec(context.Background(),
		`INSERT INTO projects (user_id, name, color) VALUES ($1, $2, $3)`,
		userID, "Project 1", "#FF0000")
	_, _ = db.Pool.Exec(context.Background(),
		`INSERT INTO projects (user_id, name, color) VALUES ($1, $2, $3)`,
		userID, "Project 2", "#00FF00")

	t.Run("list projects", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/projects", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var projects []models.Project
		err := json.Unmarshal(w.Body.Bytes(), &projects)
		assert.NoError(t, err)
		assert.Len(t, projects, 2)
	})
}

func TestProjectUpdate(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	handler := NewProjectHandler(db)
	router := setupProjectRouter(handler, userID)

	var projectID int
	err := db.Pool.QueryRow(context.Background(),
		`INSERT INTO projects (user_id, name, color) VALUES ($1, $2, $3) RETURNING id`,
		userID, "Original Project", "#FF0000").Scan(&projectID)
	assert.NoError(t, err)

	t.Run("update project successfully", func(t *testing.T) {
		reqBody := models.UpdateProjectRequest{
			Name:  stringPtr("Updated Project"),
			Color: stringPtr("#00FF00"),
		}

		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/projects/%d", projectID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var project models.Project
		err := json.Unmarshal(w.Body.Bytes(), &project)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Project", project.Name)
		assert.Equal(t, "#00FF00", project.Color)
	})
}

func TestProjectDelete(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userID := createTestUser(t, db)
	handler := NewProjectHandler(db)
	router := setupProjectRouter(handler, userID)

	var projectID int
	err := db.Pool.QueryRow(context.Background(),
		`INSERT INTO projects (user_id, name, color) VALUES ($1, $2, $3) RETURNING id`,
		userID, "Project to Delete", "#FF0000").Scan(&projectID)
	assert.NoError(t, err)

	t.Run("delete project successfully", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", fmt.Sprintf("/projects/%d", projectID), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify it's deleted
		var count int
		db.Pool.QueryRow(context.Background(), "SELECT COUNT(*) FROM projects WHERE id = $1", projectID).Scan(&count)
		assert.Equal(t, 0, count)
	})
}

func stringPtr(s string) *string {
	return &s
}
