package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/care4u/services/booking/internal/config"
	"github.com/care4u/services/booking/internal/events"
	"github.com/care4u/services/booking/internal/repository/postgres"
	transport "github.com/care4u/services/booking/internal/transport/http"
	"github.com/care4u/services/booking/internal/usecase"
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

	// NATS Event Producer
	producer, err := events.NewEventProducer(cfg.NatsURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to NATS, events will not be published: %v", err)
	} else if producer != nil {
		defer producer.Close()
	}

	// Dependency Injection
	bookingRepo := postgres.NewBookingRepository(db)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepo, producer)
	bookingHandler := transport.NewBookingHandler(bookingUsecase)

	// Router Setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up", "service": "booking"})
	})

	bookingHandler.RegisterRoutes(r)

	// Start Server
	log.Printf("Booking service starting on port %s", cfg.Port)
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
