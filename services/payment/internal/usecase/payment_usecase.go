package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/care4u/services/payment/internal/config"
	"github.com/care4u/services/payment/internal/domain"
	"github.com/google/uuid"
)

type paymentUsecase struct {
	repo   domain.PaymentRepository
	config *config.Config
}

// NewPaymentUsecase creates a new payment usecase
func NewPaymentUsecase(repo domain.PaymentRepository, cfg *config.Config) domain.PaymentUsecase {
	return &paymentUsecase{repo: repo, config: cfg}
}

func (u *paymentUsecase) InitiatePayment(ctx context.Context, req *domain.InitiatePaymentRequest) (*domain.InitiatePaymentResponse, error) {
	currency := req.Currency
	if currency == "" {
		currency = "INR"
	}

	// Create Razorpay order (mock for now)
	orderID := "order_" + uuid.New().String()[:8]

	now := time.Now()
	paymentID := uuid.New().String()
	
	payment := &domain.Payment{
		ID:        paymentID,
		BookingID: req.BookingID,
		UserID:    req.UserID,
		Amount:    req.Amount,
		Currency:  currency,
		Status:    domain.PaymentStatusCreated,
		GatewayID: orderID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.repo.Create(ctx, payment); err != nil {
		return nil, fmt.Errorf("failed to create payment record: %w", err)
	}

	// In production, you would call Razorpay API here:
	// client := razorpay.NewClient(u.config.RazorpayKeyID, u.config.RazorpayKeySecret)
	// order, err := client.Order.Create(...)

	return &domain.InitiatePaymentResponse{
		PaymentID: paymentID,
		OrderID:   orderID,
		Amount:    req.Amount,
		Currency:  currency,
		KeyID:     u.config.RazorpayKeyID,
	}, nil
}

func (u *paymentUsecase) HandleWebhook(ctx context.Context, payload *domain.WebhookPayload) error {
	// In production, verify signature first
	
	switch payload.Event {
	case "payment.captured":
		return u.handlePaymentCaptured(ctx, payload)
	case "payment.failed":
		return u.handlePaymentFailed(ctx, payload)
	default:
		// Ignore other events
		return nil
	}
}

func (u *paymentUsecase) handlePaymentCaptured(ctx context.Context, payload *domain.WebhookPayload) error {
	paymentData, ok := payload.Payload["payment"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment data in webhook")
	}
	
	entity, ok := paymentData["entity"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid entity in payment data")
	}

	orderID, _ := entity["order_id"].(string)
	paymentID, _ := entity["id"].(string)

	payment, err := u.repo.GetByGatewayID(ctx, orderID)
	if err != nil || payment == nil {
		return fmt.Errorf("payment not found for order: %s", orderID)
	}

	return u.repo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusCaptured, paymentID, "")
}

func (u *paymentUsecase) handlePaymentFailed(ctx context.Context, payload *domain.WebhookPayload) error {
	paymentData, ok := payload.Payload["payment"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid payment data in webhook")
	}
	
	entity, ok := paymentData["entity"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid entity in payment data")
	}

	orderID, _ := entity["order_id"].(string)
	reason, _ := entity["error_description"].(string)

	payment, err := u.repo.GetByGatewayID(ctx, orderID)
	if err != nil || payment == nil {
		return fmt.Errorf("payment not found for order: %s", orderID)
	}

	return u.repo.UpdateStatus(ctx, payment.ID, domain.PaymentStatusFailed, "", reason)
}

func (u *paymentUsecase) GetPaymentStatus(ctx context.Context, id string) (*domain.Payment, error) {
	payment, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}
	if payment == nil {
		return nil, fmt.Errorf("payment not found")
	}
	return payment, nil
}
