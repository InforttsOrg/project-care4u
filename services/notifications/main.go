package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/care4u/services/notifications/internal/config"
	"github.com/care4u/services/notifications/internal/domain"
	"github.com/care4u/services/notifications/internal/repository/memory"
	transport "github.com/care4u/services/notifications/internal/transport/http"
	"github.com/care4u/services/notifications/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

type BookingEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
	Data      struct {
		BookingID   string    `json:"booking_id"`
		PatientID   string    `json:"patient_id"`
		ProviderID  string    `json:"provider_id"`
		ServiceType string    `json:"service_type"`
		StartTime   time.Time `json:"start_time"`
		EndTime     time.Time `json:"end_time"`
		Status      string    `json:"status"`
		Notes       string    `json:"notes,omitempty"`
	} `json:"data"`
}

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

	// NATS Consumer for event-driven notifications
	go startNATSConsumer(cfg.NatsURL, notificationUsecase)

	// Start Server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("Notification service starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	log.Println("Server exited")
}

func startNATSConsumer(natsURL string, usecase domain.NotificationUsecase) {
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Printf("Failed to connect to NATS for consumer: %v", err)
		return
	}
	defer nc.Drain()

	js, err := nc.JetStream()
	if err != nil {
		log.Printf("Failed to get JetStream context: %v", err)
		return
	}

	// Create stream if not exists
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "NOTIFICATIONS",
		Subjects: []string{"booking.*"},
	})
	if err != nil {
		log.Printf("Stream 'NOTIFICATIONS' might already exist: %v", err)
	}

	// Create durable consumer
	_, err = js.AddConsumer("NOTIFICATIONS", &nats.ConsumerConfig{
		Durable:       "notifications-service",
		FilterSubject: "booking.*",
		AckPolicy:     nats.AckExplicitPolicy,
	})
	if err != nil {
		log.Printf("Consumer might already exist: %v", err)
	}

	// Subscribe to booking events
	sub, err := js.PullSubscribe("booking.*", "notifications-service", nats.BindStream("NOTIFICATIONS"))
	if err != nil {
		log.Printf("Failed to create pull subscription: %v", err)
		return
	}

	log.Println("NATS consumer started for booking events")

	for {
		msgs, err := sub.Fetch(10, nats.MaxWait(5*time.Second))
		if err != nil {
			if err == nats.ErrTimeout {
				continue
			}
			log.Printf("Error fetching messages: %v", err)
			continue
		}

		for _, msg := range msgs {
			var event BookingEvent
			if err := json.Unmarshal(msg.Data, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				msg.Nak()
				continue
			}

			// Process event based on type
			processBookingEvent(event, usecase)
			msg.Ack()
		}
	}
}

func processBookingEvent(event BookingEvent, usecase domain.NotificationUsecase) {
	ctx := context.Background()

	var title, body string
	var userIDs []string
	data := map[string]string{
		"booking_id":   event.Data.BookingID,
		"service_type": event.Data.ServiceType,
		"start_time":   event.Data.StartTime.Format(time.RFC3339),
		"end_time":     event.Data.EndTime.Format(time.RFC3339),
	}

	switch event.EventType {
	case "booking.created":
		title = "New Appointment Booked"
		body = "Your appointment has been booked successfully"
		userIDs = []string{event.Data.PatientID, event.Data.ProviderID}
	case "booking.confirmed":
		title = "Appointment Confirmed"
		body = "Your appointment has been confirmed"
		userIDs = []string{event.Data.PatientID, event.Data.ProviderID}
	case "booking.cancelled":
		title = "Appointment Cancelled"
		body = "Your appointment has been cancelled"
		userIDs = []string{event.Data.PatientID, event.Data.ProviderID}
	default:
		log.Printf("Unknown event type: %s", event.EventType)
		return
	}

	// Send notifications to both patient and provider
	for _, userID := range userIDs {
		req := &domain.SendNotificationRequest{
			UserID: userID,
			Type:   domain.TypePush,
			Title:  title,
			Body:   body,
			Data:   data,
		}
		if _, err := usecase.SendNotification(ctx, req); err != nil {
			log.Printf("Failed to send notification to %s: %v", userID, err)
		}
	}
}