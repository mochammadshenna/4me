package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/config"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/handlers"
	"github.com/mochammadshenna/4me-backend/internal/middleware"
	"github.com/mochammadshenna/4me-backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type E2ETestSuite struct {
	suite.Suite
	router *gin.Engine
	db     *database.Database
	token  string
	userID int
}

func (suite *E2ETestSuite) SetupSuite() {
	gin.SetMode(gin.TestMode)

	cfg := &config.Config{
		DatabaseURL: "postgres://localhost:5432/4me_todos_test?sslmode=disable",
		JWTSecret:   "test-secret-key",
		SupabaseURL: "http://localhost:54321",
		SupabaseKey: "test-key",
		FrontendURL: "http://localhost:5173",
	}

	var err error
	suite.db, err = database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		suite.T().Skipf("Skipping E2E tests: database not available - %v", err)
		return
	}

	// Run migrations
	err = suite.db.Migrate()
	assert.NoError(suite.T(), err)

	// Setup router
	suite.router = gin.New()
	suite.router.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	authHandler := handlers.NewAuthHandler(suite.db, cfg)
	projectHandler := handlers.NewProjectHandler(suite.db)
	boardHandler := handlers.NewBoardHandler(suite.db)
	taskHandler := handlers.NewTaskHandler(suite.db)
	labelHandler := handlers.NewLabelHandler(suite.db)
	commentHandler := handlers.NewCommentHandler(suite.db)

	api := suite.router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/auth/me", authHandler.Me)
		protected.POST("/projects", projectHandler.Create)
		protected.GET("/projects", projectHandler.List)
		protected.GET("/projects/:id", projectHandler.Get)
		protected.PUT("/projects/:id", projectHandler.Update)
		protected.DELETE("/projects/:id", projectHandler.Delete)

		protected.POST("/projects/:id/boards", boardHandler.Create)
		protected.GET("/projects/:id/boards", boardHandler.List)
		protected.PUT("/boards/:id", boardHandler.Update)
		protected.DELETE("/boards/:id", boardHandler.Delete)

		protected.POST("/boards/:id/tasks", taskHandler.Create)
		protected.GET("/tasks/:id", taskHandler.Get)
		protected.PUT("/tasks/:id", taskHandler.Update)
		protected.PATCH("/tasks/:id/move", taskHandler.Move)
		protected.DELETE("/tasks/:id", taskHandler.Delete)
		protected.GET("/tasks/:id/history", taskHandler.GetHistory)

		protected.POST("/projects/:id/labels", labelHandler.Create)
		protected.GET("/projects/:id/labels", labelHandler.List)
		protected.PUT("/labels/:id", labelHandler.Update)
		protected.DELETE("/labels/:id", labelHandler.Delete)

		protected.POST("/tasks/:id/comments", commentHandler.Create)
		protected.GET("/tasks/:id/comments", commentHandler.List)
		protected.PUT("/comments/:id", commentHandler.Update)
		protected.DELETE("/comments/:id", commentHandler.Delete)
	}
}

func (suite *E2ETestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

func (suite *E2ETestSuite) SetupTest() {
	// Clean database before each test
	suite.db.Pool.Exec(context.Background(), "TRUNCATE users, projects, boards, tasks, labels, task_labels, comments, attachments, task_history CASCADE")
}

func (suite *E2ETestSuite) TestCompleteWorkflow() {
	// 1. Register a user
	registerReq := models.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var authResp models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &authResp)
	suite.token = authResp.Token
	suite.userID = authResp.User.ID

	// 2. Create a project
	projectReq := models.CreateProjectRequest{
		Name:        "My Todo Project",
		Description: stringPtr("A test project"),
		Color:       "#3B82F6",
	}

	body, _ = json.Marshal(projectReq)
	req = httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var project models.Project
	json.Unmarshal(w.Body.Bytes(), &project)

	// 3. Create boards
	boardReq := models.CreateBoardRequest{
		Name:     "To Do",
		Position: 0,
	}

	body, _ = json.Marshal(boardReq)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/projects/%d/boards", project.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var board models.Board
	json.Unmarshal(w.Body.Bytes(), &board)

	// 4. Create a label
	labelReq := models.CreateLabelRequest{
		Name:  "Bug",
		Color: "#EF4444",
	}

	body, _ = json.Marshal(labelReq)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/projects/%d/labels", project.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var label models.Label
	json.Unmarshal(w.Body.Bytes(), &label)

	// 5. Create a task
	taskReq := models.CreateTaskRequest{
		Title:    "Fix login bug",
		Priority: "high",
		LabelIDs: []int{label.ID},
	}

	body, _ = json.Marshal(taskReq)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/boards/%d/tasks", board.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var task models.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	// 6. Add a comment to the task
	commentReq := models.CreateCommentRequest{
		Content: "This is a critical bug",
	}

	body, _ = json.Marshal(commentReq)
	req = httptest.NewRequest("POST", fmt.Sprintf("/api/tasks/%d/comments", task.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	// 7. Get task with comments
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tasks/%d/comments", task.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var comments []models.Comment
	json.Unmarshal(w.Body.Bytes(), &comments)
	assert.Len(suite.T(), comments, 1)

	// 8. Get task history
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/tasks/%d/history", task.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var history []models.TaskHistory
	json.Unmarshal(w.Body.Bytes(), &history)
	assert.NotEmpty(suite.T(), history)

	// 9. Update task
	updateReq := models.UpdateTaskRequest{
		Title: stringPtr("Fix critical login bug"),
	}

	body, _ = json.Marshal(updateReq)
	req = httptest.NewRequest("PUT", fmt.Sprintf("/api/tasks/%d", task.ID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	// 10. Delete task
	req = httptest.NewRequest("DELETE", fmt.Sprintf("/api/tasks/%d", task.ID), nil)
	req.Header.Set("Authorization", "Bearer "+suite.token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *E2ETestSuite) TestUnauthorizedAccess() {
	// Try to access protected endpoint without token
	req := httptest.NewRequest("GET", "/api/projects", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *E2ETestSuite) TestProjectOwnership() {
	// Register first user
	registerReq := models.RegisterRequest{
		Username: "user1",
		Email:    "user1@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(registerReq)
	req := httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var authResp1 models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &authResp1)

	// Register second user
	registerReq2 := models.RegisterRequest{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "password123",
	}

	body, _ = json.Marshal(registerReq2)
	req = httptest.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var authResp2 models.AuthResponse
	json.Unmarshal(w.Body.Bytes(), &authResp2)

	// User 1 creates a project
	projectReq := models.CreateProjectRequest{
		Name:  "User 1 Project",
		Color: "#3B82F6",
	}

	body, _ = json.Marshal(projectReq)
	req = httptest.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+authResp1.Token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var project models.Project
	json.Unmarshal(w.Body.Bytes(), &project)

	// User 2 tries to access User 1's project
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/projects/%d", project.ID), nil)
	req.Header.Set("Authorization", "Bearer "+authResp2.Token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)

	// User 1 should be able to access their own project
	req = httptest.NewRequest("GET", fmt.Sprintf("/api/projects/%d", project.ID), nil)
	req.Header.Set("Authorization", "Bearer "+authResp1.Token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func TestE2ETestSuite(t *testing.T) {
	suite.Run(t, new(E2ETestSuite))
}

func stringPtr(s string) *string {
	return &s
}
