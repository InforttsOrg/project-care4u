package usecase

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/care4u/services/booking/internal/domain"
	"github.com/care4u/services/booking/internal/events"
	"github.com/google/uuid"
)

type bookingUsecase struct {
	repo     domain.BookingRepository
	producer events.EventProducer
}

// NewBookingUsecase creates a new booking usecase
func NewBookingUsecase(repo domain.BookingRepository, producer events.EventProducer) domain.BookingUsecase {
	return &bookingUsecase{repo: repo, producer: producer}
}

func (u *bookingUsecase) CreateBooking(ctx context.Context, req *domain.CreateBookingRequest) (*domain.Booking, error) {
	// Validate time slot
	if req.StartTime.Before(time.Now()) {
		return nil, fmt.Errorf("cannot book appointments in the past")
	}
	if req.EndTime.Before(req.StartTime) || req.EndTime.Equal(req.StartTime) {
		return nil, fmt.Errorf("end time must be after start time")
	}

	// Check availability
	available, err := u.repo.CheckAvailability(ctx, req.ProviderID, req.StartTime, req.EndTime)
	if err != nil {
		return nil, fmt.Errorf("failed to check availability: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("time slot is not available")
	}

	now := time.Now()
	booking := &domain.Booking{
		ID:          uuid.New().String(),
		PatientID:   req.PatientID,
		ProviderID:  req.ProviderID,
		ServiceType: req.ServiceType,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      domain.StatusPending,
		Notes:       req.Notes,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.repo.Create(ctx, booking); err != nil {
		return nil, fmt.Errorf("failed to create booking: %w", err)
	}

	// Publish event to NATS for notifications
	event := events.NewBookingCreatedEvent(booking)
	if u.producer != nil {
		if err := u.producer.PublishBookingCreated(event); err != nil {
			log.Printf("Failed to publish booking created event: %v", err)
		}
	}

	return booking, nil
}

func (u *bookingUsecase) GetBooking(ctx context.Context, id string) (*domain.Booking, error) {
	booking, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return nil, fmt.Errorf("booking not found")
	}
	return booking, nil
}

func (u *bookingUsecase) GetPatientBookings(ctx context.Context, patientID string, limit, offset int) ([]*domain.Booking, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.GetByPatientID(ctx, patientID, limit, offset)
}

func (u *bookingUsecase) GetProviderBookings(ctx context.Context, providerID string, limit, offset int) ([]*domain.Booking, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.GetByProviderID(ctx, providerID, limit, offset)
}

func (u *bookingUsecase) CancelBooking(ctx context.Context, id, userID string) error {
	booking, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	// Only patient or provider can cancel
	if booking.PatientID != userID && booking.ProviderID != userID {
		return fmt.Errorf("unauthorized to cancel this booking")
	}

	// Can't cancel completed or already cancelled bookings
	if booking.Status == domain.StatusCompleted || booking.Status == domain.StatusCancelled {
		return fmt.Errorf("cannot cancel booking with status: %s", booking.Status)
	}

	if err := u.repo.UpdateStatus(ctx, id, domain.StatusCancelled); err != nil {
		return fmt.Errorf("failed to cancel booking: %w", err)
	}

	// Fetch updated booking for event
	booking, err = u.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Failed to fetch booking for cancellation event: %v", err)
		return nil
	}

	// Publish cancellation event
	event := events.NewBookingCancelledEvent(booking)
	if u.producer != nil {
		if err := u.producer.PublishBookingCancelled(event); err != nil {
			log.Printf("Failed to publish booking cancelled event: %v", err)
		}
	}

	return nil
}

func (u *bookingUsecase) ConfirmBooking(ctx context.Context, id string) error {
	booking, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get booking: %w", err)
	}
	if booking == nil {
		return fmt.Errorf("booking not found")
	}

	if booking.Status != domain.StatusPending {
		return fmt.Errorf("can only confirm pending bookings")
	}

	if err := u.repo.UpdateStatus(ctx, id, domain.StatusConfirmed); err != nil {
		return fmt.Errorf("failed to confirm booking: %w", err)
	}

	// Fetch updated booking for event
	booking, err = u.repo.GetByID(ctx, id)
	if err != nil {
		log.Printf("Failed to fetch booking for confirmation event: %v", err)
		return nil
	}

	// Publish confirmation event
	event := events.NewBookingConfirmedEvent(booking)
	if u.producer != nil {
		if err := u.producer.PublishBookingConfirmed(event); err != nil {
			log.Printf("Failed to publish booking confirmed event: %v", err)
		}
	}

	return nil
}
