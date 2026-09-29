package domain

import (
	"context"
	"time"
)

// NotificationType represents the type of notification
type NotificationType string

const (
	TypePush  NotificationType = "push"
	TypeSMS   NotificationType = "sms"
	TypeEmail NotificationType = "email"
	TypeInApp NotificationType = "in_app"
)

// NotificationStatus represents the delivery status
type NotificationStatus string

const (
	StatusQueued    NotificationStatus = "queued"
	StatusSent      NotificationStatus = "sent"
	StatusDelivered NotificationStatus = "delivered"
	StatusFailed    NotificationStatus = "failed"
)

// Notification represents a notification record
type Notification struct {
	ID         string             `json:"id"`
	UserID     string             `json:"user_id"`
	Type       NotificationType   `json:"type"`
	Title      string             `json:"title"`
	Body       string             `json:"body"`
	Data       map[string]string  `json:"data,omitempty"`
	Status     NotificationStatus `json:"status"`
	FailReason string             `json:"fail_reason,omitempty"`
	SentAt     *time.Time         `json:"sent_at,omitempty"`
	CreatedAt  time.Time          `json:"created_at"`
}

// SendNotificationRequest represents a notification request
type SendNotificationRequest struct {
	UserID string            `json:"user_id" binding:"required"`
	Type   NotificationType  `json:"type" binding:"required"`
	Title  string            `json:"title" binding:"required"`
	Body   string            `json:"body" binding:"required"`
	Data   map[string]string `json:"data"`
}

// NotificationRepository defines persistence operations
type NotificationRepository interface {
	Create(ctx context.Context, notification *Notification) error
	GetByID(ctx context.Context, id string) (*Notification, error)
	GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*Notification, error)
	UpdateStatus(ctx context.Context, id string, status NotificationStatus, failReason string) error
}

// NotificationUsecase defines business logic
type NotificationUsecase interface {
	SendNotification(ctx context.Context, req *SendNotificationRequest) (*Notification, error)
	GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*Notification, error)
}

// PushService defines the push notification sender
type PushService interface {
	Send(ctx context.Context, userID, title, body string, data map[string]string) error
}
