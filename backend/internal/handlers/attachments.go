package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/config"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/models"
)

type AttachmentHandler struct {
	db     *database.Database
	config *config.Config
}

func NewAttachmentHandler(db *database.Database, cfg *config.Config) *AttachmentHandler {
	return &AttachmentHandler{
		db:     db,
		config: cfg,
	}
}

func (h *AttachmentHandler) Upload(c *gin.Context) {
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

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()

	// Upload to Supabase Storage
	fileURL, err := h.uploadToSupabase(file, header, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	// Save attachment record
	var attachment models.Attachment
	err = h.db.Pool.QueryRow(context.Background(),
		`INSERT INTO attachments (task_id, filename, file_url, file_type, size) 
		 VALUES ($1, $2, $3, $4, $5) 
		 RETURNING id, task_id, filename, file_url, file_type, size, uploaded_at`,
		taskID, header.Filename, fileURL, header.Header.Get("Content-Type"), header.Size).
		Scan(&attachment.ID, &attachment.TaskID, &attachment.Filename, &attachment.FileURL, &attachment.FileType, &attachment.Size, &attachment.UploadedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save attachment"})
		return
	}

	c.JSON(http.StatusCreated, attachment)
}

func (h *AttachmentHandler) uploadToSupabase(file multipart.File, header *multipart.FileHeader, taskID int) (string, error) {
	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Generate unique filename
	filename := fmt.Sprintf("tasks/%d/%d-%s", taskID, time.Now().Unix(), header.Filename)

	// Upload to Supabase Storage
	url := fmt.Sprintf("%s/storage/v1/object/4me-attachments/%s", h.config.SupabaseURL, filename)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(fileBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+h.config.SupabaseKey)
	req.Header.Set("Content-Type", header.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase upload failed: %s", string(body))
	}

	// Return public URL
	publicURL := fmt.Sprintf("%s/storage/v1/object/public/4me-attachments/%s", h.config.SupabaseURL, filename)
	return publicURL, nil
}

func (h *AttachmentHandler) List(c *gin.Context) {
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
		`SELECT id, task_id, filename, file_url, file_type, size, uploaded_at 
		 FROM attachments WHERE task_id = $1 ORDER BY uploaded_at DESC`,
		taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attachments"})
		return
	}
	defer rows.Close()

	attachments := []models.Attachment{}
	for rows.Next() {
		var attachment models.Attachment
		if err := rows.Scan(&attachment.ID, &attachment.TaskID, &attachment.Filename, &attachment.FileURL, &attachment.FileType, &attachment.Size, &attachment.UploadedAt); err != nil {
			continue
		}
		attachments = append(attachments, attachment)
	}

	c.JSON(http.StatusOK, attachments)
}

func (h *AttachmentHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("userID")
	attachmentID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid attachment ID"})
		return
	}

	// Get attachment details and verify ownership
	var fileURL string
	var taskID int
	err = h.db.Pool.QueryRow(context.Background(),
		`SELECT a.file_url, a.task_id 
		 FROM attachments a
		 JOIN tasks t ON a.task_id = t.id 
		 JOIN boards b ON t.board_id = b.id 
		 JOIN projects p ON b.project_id = p.id 
		 WHERE a.id = $1 AND p.user_id = $2`,
		attachmentID, userID).Scan(&fileURL, &taskID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attachment not found"})
		return
	}

	// Delete from database
	result, err := h.db.Pool.Exec(context.Background(),
		"DELETE FROM attachments WHERE id = $1",
		attachmentID)

	if err != nil || result.RowsAffected() == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attachment"})
		return
	}

	// TODO: Delete from Supabase Storage
	// This would require parsing the fileURL and calling Supabase delete API

	c.JSON(http.StatusOK, gin.H{"message": "Attachment deleted successfully"})
}
