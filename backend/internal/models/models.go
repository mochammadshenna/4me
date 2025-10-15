package models

import "time"

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	GoogleID     *string   `json:"google_id,omitempty"`
	AvatarURL    *string   `json:"avatar_url,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Project struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Board struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Task struct {
	ID          int        `json:"id"`
	BoardID     int        `json:"board_id"`
	Title       string     `json:"title"`
	Description *string    `json:"description,omitempty"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  *int       `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	Position    int        `json:"position"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Labels      []Label    `json:"labels,omitempty"`
}

type Label struct {
	ID        int       `json:"id"`
	ProjectID int       `json:"project_id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
}

type Comment struct {
	ID        int       `json:"id"`
	TaskID    int       `json:"task_id"`
	UserID    int       `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	User      *User     `json:"user,omitempty"`
}

type Attachment struct {
	ID         int       `json:"id"`
	TaskID     int       `json:"task_id"`
	Filename   string    `json:"filename"`
	FileURL    string    `json:"file_url"`
	FileType   *string   `json:"file_type,omitempty"`
	Size       *int64    `json:"size,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

type TaskHistory struct {
	ID          int                    `json:"id"`
	TaskID      int                    `json:"task_id"`
	UserID      int                    `json:"user_id"`
	Action      string                 `json:"action"`
	ChangesJSON map[string]interface{} `json:"changes"`
	CreatedAt   time.Time              `json:"created_at"`
	User        *User                  `json:"user,omitempty"`
}

// Request/Response DTOs

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	User         User   `json:"user"`
}

type CreateProjectRequest struct {
	Name        string  `json:"name" binding:"required,min=1,max=255"`
	Description *string `json:"description"`
	Color       string  `json:"color"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Color       *string `json:"color"`
}

type CreateBoardRequest struct {
	Name     string `json:"name" binding:"required,min=1,max=255"`
	Position int    `json:"position"`
}

type UpdateBoardRequest struct {
	Name     *string `json:"name"`
	Position *int    `json:"position"`
}

type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=255"`
	Description *string    `json:"description"`
	Priority    string     `json:"priority"`
	AssigneeID  *int       `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	LabelIDs    []int      `json:"label_ids"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Status      *string    `json:"status"`
	Priority    *string    `json:"priority"`
	AssigneeID  *int       `json:"assignee_id"`
	DueDate     *time.Time `json:"due_date"`
	Position    *int       `json:"position"`
	LabelIDs    []int      `json:"label_ids"`
}

type MoveTaskRequest struct {
	BoardID  int `json:"board_id" binding:"required"`
	Position int `json:"position"`
}

type CreateLabelRequest struct {
	Name  string `json:"name" binding:"required,min=1,max=100"`
	Color string `json:"color"`
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}

type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required,min=1"`
}
