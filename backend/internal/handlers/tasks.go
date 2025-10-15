package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type TaskHandler struct {
	db *database.Database
}

func NewTaskHandler(db *database.Database) *TaskHandler {
	return &TaskHandler{db: db}
}

func (h *TaskHandler) Create(c *gin.Context) {
	userID, _ := c.Get("userID")
	boardID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board ID"})
		return
	}

	// Verify board ownership
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

	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Priority == "" {
		req.Priority = "medium"
	}

	ctx := context.Background()
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	var task models.Task
	err = tx.QueryRow(ctx,
		`INSERT INTO tasks (board_id, title, description, priority, assignee_id, due_date, status) 
		 VALUES ($1, $2, $3, $4, $5, $6, 'todo') 
		 RETURNING id, board_id, title, description, status, priority, assignee_id, due_date, position, created_at, updated_at`,
		boardID, req.Title, req.Description, req.Priority, req.AssigneeID, req.DueDate).
		Scan(&task.ID, &task.BoardID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.AssigneeID, &task.DueDate, &task.Position, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	// Add labels
	if len(req.LabelIDs) > 0 {
		for _, labelID := range req.LabelIDs {
			_, err = tx.Exec(ctx, "INSERT INTO task_labels (task_id, label_id) VALUES ($1, $2)", task.ID, labelID)
			if err != nil {
				continue
			}
		}
	}

	// Add to history
	changes := map[string]interface{}{
		"action": "created",
		"title":  req.Title,
	}
	changesJSON, _ := json.Marshal(changes)
	_, _ = tx.Exec(ctx,
		"INSERT INTO task_history (task_id, user_id, action, changes_json) VALUES ($1, $2, $3, $4)",
		task.ID, userID, "created", changesJSON)

	if err = tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) Get(c *gin.Context) {
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

	var task models.Task
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT id, board_id, title, description, status, priority, assignee_id, due_date, position, created_at, updated_at 
		 FROM tasks WHERE id = $1`,
		taskID).
		Scan(&task.ID, &task.BoardID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.AssigneeID, &task.DueDate, &task.Position, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Get labels
	rows, err := h.db.Pool.Query(context.Background(),
		`SELECT l.id, l.project_id, l.name, l.color, l.created_at 
		 FROM labels l 
		 JOIN task_labels tl ON l.id = tl.label_id 
		 WHERE tl.task_id = $1`,
		taskID)
	if err == nil {
		defer rows.Close()
		task.Labels = []models.Label{}
		for rows.Next() {
			var label models.Label
			if err := rows.Scan(&label.ID, &label.ProjectID, &label.Name, &label.Color, &label.CreatedAt); err == nil {
				task.Labels = append(task.Labels, label)
			}
		}
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c *gin.Context) {
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

	var req models.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	query := "UPDATE tasks SET updated_at = NOW()"
	args := []interface{}{}
	argCount := 1
	changes := map[string]interface{}{}

	if req.Title != nil {
		query += ", title = $" + strconv.Itoa(argCount)
		args = append(args, *req.Title)
		changes["title"] = *req.Title
		argCount++
	}
	if req.Description != nil {
		query += ", description = $" + strconv.Itoa(argCount)
		args = append(args, *req.Description)
		changes["description"] = *req.Description
		argCount++
	}
	if req.Status != nil {
		query += ", status = $" + strconv.Itoa(argCount)
		args = append(args, *req.Status)
		changes["status"] = *req.Status
		argCount++
	}
	if req.Priority != nil {
		query += ", priority = $" + strconv.Itoa(argCount)
		args = append(args, *req.Priority)
		changes["priority"] = *req.Priority
		argCount++
	}
	if req.AssigneeID != nil {
		query += ", assignee_id = $" + strconv.Itoa(argCount)
		args = append(args, *req.AssigneeID)
		changes["assignee_id"] = *req.AssigneeID
		argCount++
	}
	if req.DueDate != nil {
		query += ", due_date = $" + strconv.Itoa(argCount)
		args = append(args, *req.DueDate)
		changes["due_date"] = *req.DueDate
		argCount++
	}
	if req.Position != nil {
		query += ", position = $" + strconv.Itoa(argCount)
		args = append(args, *req.Position)
		argCount++
	}

	query += " WHERE id = $" + strconv.Itoa(argCount)
	args = append(args, taskID)
	query += " RETURNING id, board_id, title, description, status, priority, assignee_id, due_date, position, created_at, updated_at"

	var task models.Task
	err = tx.QueryRow(ctx, query, args...).
		Scan(&task.ID, &task.BoardID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.AssigneeID, &task.DueDate, &task.Position, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Update labels if provided
	if req.LabelIDs != nil {
		_, _ = tx.Exec(ctx, "DELETE FROM task_labels WHERE task_id = $1", taskID)
		for _, labelID := range req.LabelIDs {
			_, _ = tx.Exec(ctx, "INSERT INTO task_labels (task_id, label_id) VALUES ($1, $2)", taskID, labelID)
		}
		changes["labels"] = req.LabelIDs
	}

	// Add to history
	if len(changes) > 0 {
		changesJSON, _ := json.Marshal(changes)
		_, _ = tx.Exec(ctx,
			"INSERT INTO task_history (task_id, user_id, action, changes_json) VALUES ($1, $2, $3, $4)",
			taskID, userID, "updated", changesJSON)
	}

	if err = tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Move(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	var req models.MoveTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	tx, err := h.db.Pool.Begin(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}
	defer tx.Rollback(ctx)

	// Verify ownership and update
	var task models.Task
	err = tx.QueryRow(ctx,
		`UPDATE tasks SET board_id = $1, position = $2, updated_at = NOW() 
		 WHERE id = $3 AND board_id IN (
			 SELECT b.id FROM boards b 
			 JOIN projects p ON b.project_id = p.id 
			 WHERE p.user_id = $4
		 )
		 RETURNING id, board_id, title, description, status, priority, assignee_id, due_date, position, created_at, updated_at`,
		req.BoardID, req.Position, taskID, userID).
		Scan(&task.ID, &task.BoardID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.AssigneeID, &task.DueDate, &task.Position, &task.CreatedAt, &task.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or unauthorized"})
		return
	}

	// Add to history
	changes := map[string]interface{}{
		"board_id": req.BoardID,
		"position": req.Position,
	}
	changesJSON, _ := json.Marshal(changes)
	_, _ = tx.Exec(ctx,
		"INSERT INTO task_history (task_id, user_id, action, changes_json) VALUES ($1, $2, $3, $4)",
		taskID, userID, "moved", changesJSON)

	if err = tx.Commit(ctx); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	taskID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	result, err := h.db.Pool.Exec(context.Background(),
		`DELETE FROM tasks 
		 WHERE id = $1 AND board_id IN (
			 SELECT b.id FROM boards b 
			 JOIN projects p ON b.project_id = p.id 
			 WHERE p.user_id = $2
		 )`,
		taskID, userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	if result.RowsAffected() == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) GetHistory(c *gin.Context) {
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
		`SELECT h.id, h.task_id, h.user_id, h.action, h.changes_json, h.created_at,
		        u.id, u.username, u.email, u.avatar_url
		 FROM task_history h
		 JOIN users u ON h.user_id = u.id
		 WHERE h.task_id = $1
		 ORDER BY h.created_at DESC`,
		taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch history"})
		return
	}
	defer rows.Close()

	history := []models.TaskHistory{}
	for rows.Next() {
		var h models.TaskHistory
		var changesJSON []byte
		var user models.User

		err := rows.Scan(&h.ID, &h.TaskID, &h.UserID, &h.Action, &changesJSON, &h.CreatedAt,
			&user.ID, &user.Username, &user.Email, &user.AvatarURL)
		if err != nil {
			continue
		}

		if len(changesJSON) > 0 {
			json.Unmarshal(changesJSON, &h.ChangesJSON)
		}

		h.User = &user
		history = append(history, h)
	}

	c.JSON(http.StatusOK, history)
}

