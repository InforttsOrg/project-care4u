package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/care4u/services/auth/internal/config"
	"github.com/care4u/services/auth/internal/event"
	"github.com/care4u/services/auth/internal/repository/postgres"
	transport "github.com/care4u/services/auth/internal/transport/http"
	"github.com/care4u/services/auth/internal/usecase"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.Load()

	// DB Connection
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	// Run Migrations (Ideally use a migration tool)
	runMigrations(db)

	// NATS Connection
	eventProducer, err := event.NewEventProducer(cfg.NatsURL)
	if err != nil {
		// For MVP, maybe we don't crash, or maybe we do. usage depends.
		// The producer stays nil and the usecase skips the user.verified event
		// (it already logs publish failures as warnings), so OTP logins keep
		// working while NATS is down instead of panicking mid-request.
		log.Printf("Warning: Failed to connect to NATS (user.verified events disabled): %v", err)
	}

	// Dependency Injection
	userRepo := postgres.NewUserRepository(db)
	authUsecase := usecase.NewAuthUsecase(userRepo, eventProducer, cfg)
	authHandler := transport.NewAuthHandler(authUsecase)

	// Router Setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up"})
	})

	authHandler.RegisterRoutes(r)

	// Start Server
	log.Printf("Auth service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func runMigrations(db *sql.DB) {
	// Simple migration runner for MVP
	migration, err := os.ReadFile("../../migrations/001_init.sql")
	if err != nil {
		// Try relative path from docker container root usually /app
		// But in docker we copy migrations folder.
		// Let's assume we copy migrations to /app/migrations
		migration, err = os.ReadFile("migrations/001_init.sql")
		if err != nil {
			log.Printf("Warning: Could not read migration file: %v", err)
			return
		}
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}
	fmt.Println("Migrations executed successfully")
}
