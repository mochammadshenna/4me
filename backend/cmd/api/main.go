package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/mochammadshenna/4me-backend/internal/config"
	"github.com/mochammadshenna/4me-backend/internal/database"
	"github.com/mochammadshenna/4me-backend/internal/handlers"
	"github.com/mochammadshenna/4me-backend/internal/middleware"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Run migrations
	if err := db.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORSMiddleware(cfg.FrontendURL))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db, cfg)
	projectHandler := handlers.NewProjectHandler(db)
	boardHandler := handlers.NewBoardHandler(db)
	taskHandler := handlers.NewTaskHandler(db)
	labelHandler := handlers.NewLabelHandler(db)
	commentHandler := handlers.NewCommentHandler(db)
	attachmentHandler := handlers.NewAttachmentHandler(db, cfg)

	// Public routes
	api := router.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/google", authHandler.GoogleLogin)
			auth.GET("/google/callback", authHandler.GoogleCallback)
		}
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// User routes
		protected.GET("/auth/me", authHandler.Me)

		// Project routes
		protected.POST("/projects", projectHandler.Create)
		protected.GET("/projects", projectHandler.List)
		protected.GET("/projects/:id", projectHandler.Get)
		protected.PUT("/projects/:id", projectHandler.Update)
		protected.DELETE("/projects/:id", projectHandler.Delete)

		// Board routes (under projects)
		protected.POST("/projects/:id/boards", boardHandler.Create)
		protected.GET("/projects/:id/boards", boardHandler.List)
		protected.PUT("/boards/:id", boardHandler.Update)
		protected.DELETE("/boards/:id", boardHandler.Delete)

		// Task routes
		protected.POST("/boards/:id/tasks", taskHandler.Create)
		protected.GET("/tasks/:id", taskHandler.Get)
		protected.PUT("/tasks/:id", taskHandler.Update)
		protected.PATCH("/tasks/:id/move", taskHandler.Move)
		protected.DELETE("/tasks/:id", taskHandler.Delete)
		protected.GET("/tasks/:id/history", taskHandler.GetHistory)

		// Label routes
		protected.POST("/projects/:id/labels", labelHandler.Create)
		protected.GET("/projects/:id/labels", labelHandler.List)
		protected.PUT("/labels/:id", labelHandler.Update)
		protected.DELETE("/labels/:id", labelHandler.Delete)

		// Comment routes
		protected.POST("/tasks/:id/comments", commentHandler.Create)
		protected.GET("/tasks/:id/comments", commentHandler.List)
		protected.PUT("/comments/:id", commentHandler.Update)
		protected.DELETE("/comments/:id", commentHandler.Delete)

		// Attachment routes
		protected.POST("/tasks/:id/attachments", attachmentHandler.Upload)
		protected.GET("/tasks/:id/attachments", attachmentHandler.List)
		protected.DELETE("/attachments/:id", attachmentHandler.Delete)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Start server
	log.Printf("Server starting on port %s...", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

