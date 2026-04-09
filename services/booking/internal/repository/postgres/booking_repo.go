package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/care4u/services/booking/internal/domain"
	_ "github.com/lib/pq"
)

type bookingRepository struct {
	db *sql.DB
}

// NewBookingRepository creates a new PostgreSQL booking repository
func NewBookingRepository(db *sql.DB) domain.BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	query := `
		INSERT INTO bookings (id, patient_id, provider_id, service_type, start_time, end_time, status, notes, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.ExecContext(ctx, query,
		booking.ID,
		booking.PatientID,
		booking.ProviderID,
		booking.ServiceType,
		booking.StartTime,
		booking.EndTime,
		booking.Status,
		booking.Notes,
		booking.CreatedAt,
		booking.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create booking: %w", err)
	}
	return nil
}

func (r *bookingRepository) GetByID(ctx context.Context, id string) (*domain.Booking, error) {
	query := `
		SELECT id, patient_id, provider_id, service_type, start_time, end_time, status, notes, created_at, updated_at
		FROM bookings WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var booking domain.Booking
	err := row.Scan(
		&booking.ID,
		&booking.PatientID,
		&booking.ProviderID,
		&booking.ServiceType,
		&booking.StartTime,
		&booking.EndTime,
		&booking.Status,
		&booking.Notes,
		&booking.CreatedAt,
		&booking.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	return &booking, nil
}

func (r *bookingRepository) GetByPatientID(ctx context.Context, patientID string, limit, offset int) ([]*domain.Booking, error) {
	query := `
		SELECT id, patient_id, provider_id, service_type, start_time, end_time, status, notes, created_at, updated_at
		FROM bookings WHERE patient_id = $1
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`
	return r.queryBookings(ctx, query, patientID, limit, offset)
}

func (r *bookingRepository) GetByProviderID(ctx context.Context, providerID string, limit, offset int) ([]*domain.Booking, error) {
	query := `
		SELECT id, patient_id, provider_id, service_type, start_time, end_time, status, notes, created_at, updated_at
		FROM bookings WHERE provider_id = $1
		ORDER BY start_time DESC
		LIMIT $2 OFFSET $3
	`
	return r.queryBookings(ctx, query, providerID, limit, offset)
}

func (r *bookingRepository) queryBookings(ctx context.Context, query string, args ...interface{}) ([]*domain.Booking, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query bookings: %w", err)
	}
	defer rows.Close()

	var bookings []*domain.Booking
	for rows.Next() {
		var booking domain.Booking
		err := rows.Scan(
			&booking.ID,
			&booking.PatientID,
			&booking.ProviderID,
			&booking.ServiceType,
			&booking.StartTime,
			&booking.EndTime,
			&booking.Status,
			&booking.Notes,
			&booking.CreatedAt,
			&booking.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan booking: %w", err)
		}
		bookings = append(bookings, &booking)
	}
	return bookings, rows.Err()
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id string, status domain.BookingStatus) error {
	query := `UPDATE bookings SET status = $1, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update booking status: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("booking not found")
	}
	return nil
}

func (r *bookingRepository) CheckAvailability(ctx context.Context, providerID string, startTime, endTime time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bookings 
		WHERE provider_id = $1 
		AND status NOT IN ('cancelled', 'completed', 'no_show')
		AND (
			(start_time <= $2 AND end_time > $2) OR
			(start_time < $3 AND end_time >= $3) OR
			(start_time >= $2 AND end_time <= $3)
		)
	`
	var count int
	err := r.db.QueryRowContext(ctx, query, providerID, startTime, endTime).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}
	return count == 0, nil
}
