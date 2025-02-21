package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/s0hxm/backend/internal/api/handlers"
	"github.com/s0hxm/backend/internal/api/middleware"
	"github.com/s0hxm/backend/internal/models"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Validate required environment variables
	requiredEnvVars := []string{"DATABASE_URL", "PORT"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Error: Missing required environment variable: %s", envVar)
		}
	}

	// Set up database connection
	db, err := setupDatabase()
	if err != nil {
		log.Fatalf("❌ Database connection failed: %v", err)
	}
	log.Println("✅ Connected to database successfully!")

	// Initialize router
	router := gin.Default()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Initialize handlers
	taskHandler := handlers.NewTaskHandler(db)
	userHandler := handlers.NewUserHandler(db)

	// Set up API routes
	api := router.Group("/api")
	{
		// Task routes
		api.GET("/tasks", taskHandler.GetTasks)
		api.POST("/tasks", middleware.AuthMiddleware(), taskHandler.CreateTask)
		api.PUT("/tasks/:id", middleware.AuthMiddleware(), taskHandler.UpdateTask)
		api.DELETE("/tasks/:id", middleware.AuthMiddleware(), taskHandler.DeleteTask)

		// User routes
		api.POST("/register", userHandler.Register)
		api.POST("/login", userHandler.Login)
		api.GET("/user", middleware.AuthMiddleware(), userHandler.GetUser)
	}

	// Start the server with graceful shutdown
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	go func() {
		log.Printf("🚀 Server running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	// Graceful shutdown
	shutdownGracefully(server)
}

// setupDatabase initializes the database connection
func setupDatabase() (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate database schema
	err = db.AutoMigrate(&models.Task{}, &models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

// shutdownGracefully handles server termination
func shutdownGracefully(server *http.Server) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("\n🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("❌ Server shutdown failed: %v", err)
	}
	log.Println("✅ Server gracefully stopped")
}
