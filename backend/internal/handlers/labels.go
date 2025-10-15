package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type LabelHandler struct {
	db *database.Database
}

func NewLabelHandler(db *database.Database) *LabelHandler {
	return &LabelHandler{db: db}
}

func (h *LabelHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Verify project ownership
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)",
		projectID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var req models.CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#3B82F6"
	}

	var label models.Label
	err = h.db.Pool.QueryRow(context.Background(),
		`INSERT INTO labels (project_id, name, color) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, project_id, name, color, created_at`,
		projectID, req.Name, req.Color).
		Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create label"})
		return
	}

	c.JSON(http.StatusCreated, label)
}

func (h *LabelHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	projectID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Verify project ownership
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM projects WHERE id = $1 AND user_id = $2)",
		projectID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT id, project_id, name, color, created_at 
		 FROM labels WHERE project_id = $1 ORDER BY created_at ASC`,
		projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch labels"})
		return
	}
	defer rows.Close()

	labels := []models.Label{}
	for rows.Next() {
		var label models.Label
		if err := rows.Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt); err != nil {
			continue
		}
		labels = append(labels, label)
	}

	c.JSON(http.StatusOK, labels)
}

func (h *LabelHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	labelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	// Verify label ownership through project
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM labels l 
			JOIN projects p ON l.project_id = p.id 
			WHERE l.id = $1 AND p.user_id = $2
		)`,
		labelID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	var req models.CreateLabelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var label models.Label
	err = h.db.Pool.QueryRow(context.Background(),
		`UPDATE labels SET name = $1, color = $2 
		 WHERE id = $3 
		 RETURNING id, project_id, name, color, created_at`,
		req.Name, req.Color, labelID).
		Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update label"})
		return
	}

	c.JSON(http.StatusOK, label)
}

func (h *LabelHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	labelID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid label ID"})
		return
	}

	result, err := h.db.Pool.Exec(context.Background(),
		`DELETE FROM labels 
		 WHERE id = $1 AND project_id IN (
			 SELECT id FROM projects WHERE user_id = $2
		 )`,
		labelID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete label"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Label not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Label deleted successfully"})
}

