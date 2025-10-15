package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file from backend directory
	if err := godotenv.Load(".env"); err != nil {
		if err := godotenv.Load("../.env"); err != nil {
			log.Println("No .env file found, using environment variables")
		}
	}

	// Get database URL
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL not set")
	}

	// Connect to database
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pool.Close()

	// Hash password for demo user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	// Check if demo user already exists
	var exists bool
	err = pool.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		"demo").Scan(&exists)
	if err != nil {
		log.Fatalf("Failed to check if user exists: %v", err)
	}

	if exists {
		fmt.Println("✅ Demo user already exists")
		fmt.Println("   Username: demo")
		fmt.Println("   Password: password")
		return
	}

	// Create demo user
	var userID int64
	err = pool.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		"demo", "demo@example.com", string(passwordHash)).Scan(&userID)

	if err != nil {
		log.Fatalf("Failed to create demo user: %v", err)
	}

	fmt.Println("✅ Demo user created successfully!")
	fmt.Println("   Username: demo")
	fmt.Println("   Password: password")
	fmt.Println("   Email: demo@example.com")
	fmt.Printf("   User ID: %d\n", userID)
}
