package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/care4u/services/payment/internal/domain"
	_ "github.com/lib/pq"
)

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) domain.PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Create(ctx context.Context, payment *domain.Payment) error {
	query := `
		INSERT INTO payments (id, booking_id, user_id, amount, currency, status, gateway_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		payment.ID,
		payment.BookingID,
		payment.UserID,
		payment.Amount,
		payment.Currency,
		payment.Status,
		payment.GatewayID,
		payment.CreatedAt,
		payment.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}
	return nil
}

func (r *paymentRepository) GetByID(ctx context.Context, id string) (*domain.Payment, error) {
	return r.scanPayment(ctx, "SELECT id, booking_id, user_id, amount, currency, status, gateway_id, gateway_payment_id, failure_reason, created_at, updated_at FROM payments WHERE id = $1", id)
}

func (r *paymentRepository) GetByBookingID(ctx context.Context, bookingID string) (*domain.Payment, error) {
	return r.scanPayment(ctx, "SELECT id, booking_id, user_id, amount, currency, status, gateway_id, gateway_payment_id, failure_reason, created_at, updated_at FROM payments WHERE booking_id = $1 ORDER BY created_at DESC LIMIT 1", bookingID)
}

func (r *paymentRepository) GetByGatewayID(ctx context.Context, gatewayID string) (*domain.Payment, error) {
	return r.scanPayment(ctx, "SELECT id, booking_id, user_id, amount, currency, status, gateway_id, gateway_payment_id, failure_reason, created_at, updated_at FROM payments WHERE gateway_id = $1", gatewayID)
}

func (r *paymentRepository) scanPayment(ctx context.Context, query string, arg string) (*domain.Payment, error) {
	row := r.db.QueryRowContext(ctx, query, arg)
	
	var p domain.Payment
	var gatewayPayID, failureReason sql.NullString
	
	err := row.Scan(
		&p.ID,
		&p.BookingID,
		&p.UserID,
		&p.Amount,
		&p.Currency,
		&p.Status,
		&p.GatewayID,
		&gatewayPayID,
		&failureReason,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan payment: %w", err)
	}
	
	p.GatewayPayID = gatewayPayID.String
	p.FailureReason = failureReason.String
	
	return &p, nil
}

func (r *paymentRepository) UpdateStatus(ctx context.Context, id string, status domain.PaymentStatus, gatewayPayID, failureReason string) error {
	query := `
		UPDATE payments 
		SET status = $1, gateway_payment_id = $2, failure_reason = $3, updated_at = NOW() 
		WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query, status, gatewayPayID, failureReason, id)
	if err != nil {
		return fmt.Errorf("failed to update payment: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("payment not found")
	}
	return nil
}
