package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type CommentHandler struct {
	db *database.Database
}

func NewCommentHandler(db *database.Database) *CommentHandler {
	return &CommentHandler{db: db}
}

func (h *CommentHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Verify task ownership
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM tasks t 
			JOIN boards b ON t.board_id = b.id 
			JOIN projects p ON b.project_id = p.id 
			WHERE t.id = $1 AND p.user_id = $2
		)`,
		taskID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	var req models.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comment models.Comment
	err = h.db.Pool.QueryRow(context.Background(),
		`INSERT INTO comments (task_id, user_id, content) 
		 VALUES ($1, $2, $3) 
		 RETURNING id, task_id, user_id, content, created_at, updated_at`,
		taskID, userID, req.Content).
		Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create comment"})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *CommentHandler) List(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Verify task ownership
	var exists bool
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT EXISTS(
			SELECT 1 FROM tasks t 
			JOIN boards b ON t.board_id = b.id 
			JOIN projects p ON b.project_id = p.id 
			WHERE t.id = $1 AND p.user_id = $2
		)`,
		taskID, userID).Scan(&exists)
	if err != nil || !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT c.id, c.task_id, c.user_id, c.content, c.created_at, c.updated_at,
		        u.id, u.username, u.email, u.avatar_url
		 FROM comments c
		 JOIN users u ON c.user_id = u.id
		 WHERE c.task_id = $1
		 ORDER BY c.created_at ASC`,
		taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
		return
	}
	defer rows.Close()

	comments := []models.Comment{}
	for rows.Next() {
		var comment models.Comment
		var user models.User

		err := rows.Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt,
			&user.ID, &user.Username, &user.Email, &user.AvatarURL)
		if err != nil {
			continue
		}

		comment.User = &user
		comments = append(comments, comment)
	}

	c.JSON(http.StatusOK, comments)
}

func (h *CommentHandler) Update(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	var req models.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var comment models.Comment
	err = h.db.Pool.QueryRow(context.Background(),
		`UPDATE comments SET content = $1, updated_at = NOW() 
		 WHERE id = $2 AND user_id = $3
		 RETURNING id, task_id, user_id, content, created_at, updated_at`,
		req.Content, commentID, userID).
		Scan(&comment.ID, &comment.TaskID, &comment.UserID, &comment.Content, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *CommentHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	commentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	result, err := h.db.Pool.Exec(context.Background(),
		"DELETE FROM comments WHERE id = $1 AND user_id = $2",
		commentID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or unauthorized"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
