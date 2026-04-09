package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/care4u/services/payment/internal/config"
	"github.com/care4u/services/payment/internal/repository/postgres"
	transport "github.com/care4u/services/payment/internal/transport/http"
	"github.com/care4u/services/payment/internal/usecase"
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

	// Run Migrations
	runMigrations(db)

	// Dependency Injection
	paymentRepo := postgres.NewPaymentRepository(db)
	paymentUsecase := usecase.NewPaymentUsecase(paymentRepo, cfg)
	paymentHandler := transport.NewPaymentHandler(paymentUsecase)

	// Router Setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up", "service": "payment"})
	})

	paymentHandler.RegisterRoutes(r)

	// Start Server
	log.Printf("Payment service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

func runMigrations(db *sql.DB) {
	migration, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		log.Printf("Warning: Could not read migration file: %v", err)
		return
	}
	_, err = db.Exec(string(migration))
	if err != nil {
		log.Printf("Warning: Migration may have already been applied: %v", err)
	} else {
		log.Println("Migrations executed successfully")
	}
}
