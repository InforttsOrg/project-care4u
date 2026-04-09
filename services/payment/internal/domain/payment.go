package domain

import (
	"context"
	"time"
)

// PaymentStatus represents payment state
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusCreated   PaymentStatus = "created"
	PaymentStatusCaptured  PaymentStatus = "captured"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// Payment represents a payment record
type Payment struct {
	ID            string        `json:"id"`
	BookingID     string        `json:"booking_id"`
	UserID        string        `json:"user_id"`
	Amount        int64         `json:"amount"` // Amount in paise (100 = ₹1)
	Currency      string        `json:"currency"`
	Status        PaymentStatus `json:"status"`
	GatewayID     string        `json:"gateway_id,omitempty"` // Razorpay order ID
	GatewayPayID  string        `json:"gateway_payment_id,omitempty"`
	FailureReason string        `json:"failure_reason,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// InitiatePaymentRequest represents the input for initiating a payment
type InitiatePaymentRequest struct {
	BookingID string `json:"booking_id" binding:"required"`
	UserID    string `json:"user_id" binding:"required"`
	Amount    int64  `json:"amount" binding:"required,min=100"`
	Currency  string `json:"currency"`
}

// InitiatePaymentResponse represents the response for initiating a payment
type InitiatePaymentResponse struct {
	PaymentID  string `json:"payment_id"`
	OrderID    string `json:"order_id"`     // Razorpay order ID
	Amount     int64  `json:"amount"`
	Currency   string `json:"currency"`
	KeyID      string `json:"key_id"`       // Razorpay key for client
}

// WebhookPayload represents a Razorpay webhook event
type WebhookPayload struct {
	Entity  string                 `json:"entity"`
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
}

// PaymentRepository defines persistence operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id string) (*Payment, error)
	GetByBookingID(ctx context.Context, bookingID string) (*Payment, error)
	GetByGatewayID(ctx context.Context, gatewayID string) (*Payment, error)
	UpdateStatus(ctx context.Context, id string, status PaymentStatus, gatewayPayID, failureReason string) error
}

// PaymentUsecase defines business logic
type PaymentUsecase interface {
	InitiatePayment(ctx context.Context, req *InitiatePaymentRequest) (*InitiatePaymentResponse, error)
	HandleWebhook(ctx context.Context, payload *WebhookPayload) error
	GetPaymentStatus(ctx context.Context, id string) (*Payment, error)
}
