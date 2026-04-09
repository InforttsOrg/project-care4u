package main

import (
	"log"

	"github.com/care4u/services/notifications/internal/config"
	"github.com/care4u/services/notifications/internal/repository/memory"
	transport "github.com/care4u/services/notifications/internal/transport/http"
	"github.com/care4u/services/notifications/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	// Dependency Injection
	// Using in-memory repo for MVP. Replace with PostgreSQL/DynamoDB for production
	notificationRepo := memory.NewNotificationRepository()
	notificationUsecase := usecase.NewNotificationUsecase(notificationRepo, nil) // nil pushService uses mock
	notificationHandler := transport.NewNotificationHandler(notificationUsecase)

	// Router Setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up", "service": "notifications"})
	})

	notificationHandler.RegisterRoutes(r)

	// TODO: Add NATS consumer for event-driven notifications
	// e.g., listen for booking.created, payment.captured events

	// Start Server
	log.Printf("Notification service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
