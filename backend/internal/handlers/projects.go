package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type ProjectHandler struct {
	db *database.Database
}

func NewProjectHandler(db *database.Database) *ProjectHandler {
	return &ProjectHandler{db: db}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req models.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}

	var project models.Project
	err := h.db.Pool.QueryRow(context.Background(),
		`INSERT INTO projects (user_id, name, description, color) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, user_id, name, description, color, created_at, updated_at`,
		userID, req.Name, req.Description, req.Color).
		Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.Color, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")

	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT id, user_id, name, description, color, created_at, updated_at 
		 FROM projects WHERE user_id = $1 ORDER BY created_at DESC`,
		userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}
	defer rows.Close()

	projects := []models.Project{}
	for rows.Next() {
		var project models.Project
		if err := rows.Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.Color, &project.CreatedAt, &project.UpdatedAt); err != nil {
			continue
		}
		projects = append(projects, project)
	}

	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var project models.Project
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT id, user_id, name, description, color, created_at, updated_at 
		 FROM projects WHERE id = $1 AND user_id = $2`,
		projectID, userID).
		Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.Color, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	var req models.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check ownership
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)",
		projectID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Build dynamic update query
	query := "UPDATE projects SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		query += ", name = $" + strconv.Itoa(argCount)
		args = append(args, *req.Name)
		argCount++
	}
	if req.Description != nil {
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, *req.Description)
		argCount++
	}
	if req.Color != nil {
		query += ", color = $" + strconv.Itoa(argCount)
		args = append(args, *req.Color)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, projectID)
	argCount++

	query += " AND user_id = $" + strconv.Itoa(argCount)
	args = append(args, userID)

	query += " RETURNING id, user_id, name, description, color, created_at, updated_at"

	var project models.Project
	err = h.db.Pool.QueryRow(context.Background(), query, args...).
		Scan(&project.ID, &project.UserID, &project.Name, &project.Description, &project.Color, &project.CreatedAt, &project.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		return
	}

	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	result, err := h.db.Pool.Exec(context.Background(),
		"DELETE FROM projects WHERE id = $1 AND user_id = $2",
		projectID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

