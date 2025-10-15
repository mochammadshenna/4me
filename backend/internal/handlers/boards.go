package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type BoardHandler struct {
	db *database.Database
}

func NewBoardHandler(db *database.Database) *BoardHandler {
	return &BoardHandler{db: db}
}

func (h *BoardHandler) Create(c *gin.Context) {
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

	var req models.CreateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var board models.Board
	err = h.db.Pool.QueryRow(context.Background(),
		`INSERT INTO boards (project_id, name, position) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, project_id, name, position, created_at, updated_at`,
		projectID, req.Name, req.Position).
		Scan(&board.ID, &board.ProjectID, &board.Name, &board.Position, &board.CreatedAt, &board.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board"})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func (h *BoardHandler) List(c *gin.Context) {
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
		`SELECT id, project_id, name, position, created_at, updated_at 
		 FROM boards WHERE project_id = $1 ORDER BY position ASC`,
		projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch boards"})
		return
	}
	defer rows.Close()

	boards := []models.Board{}
	for rows.Next() {
		var board models.Board
		if err := rows.Scan(&board.ID, &board.ProjectID, &board.Name, &board.Position, &board.CreatedAt, &board.UpdatedAt); err != nil {
			continue
		}
		boards = append(boards, board)
	}

	c.JSON(http.StatusOK, boards)
}

func (h *BoardHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	// Verify board ownership through project
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM boards b 
			JOIN projects p ON b.project_id = p.id 
			WHERE b.id = $1 AND p.user_id = $2
		)`,
		boardID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	var req models.UpdateBoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "UPDATE boards SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1

	if req.Name != nil {
		query += ", name = $" + strconv.Itoa(argCount)
		args = append(args, *req.Name)
		argCount++
	}
	if req.Position != nil {
		query += ", position = $" + strconv.Itoa(argCount)
		args = append(args, *req.Position)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, boardID)

	query += " RETURNING id, project_id, name, position, created_at, updated_at"

	var board models.Board
	err = h.db.Pool.QueryRow(context.Background(), query, args...).
		Scan(&board.ID, &board.ProjectID, &board.Name, &board.Position, &board.CreatedAt, &board.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update board"})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	result, err := h.db.Pool.Exec(context.Background(),
		`DELETE FROM boards 
		 WHERE id = $1 AND project_id IN (
			 SELECT id FROM projects WHERE user_id = $2
		 )`,
		boardID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete board"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Board not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Board deleted successfully"})
}

