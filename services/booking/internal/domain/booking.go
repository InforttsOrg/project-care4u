package domain

import (
	"context"
	"time"
)

// BookingStatus represents the current state of a booking
type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
	StatusNoShow    BookingStatus = "no_show"
)

// Booking represents an appointment booking
type Booking struct {
	ID          string        `json:"id"`
	PatientID   string        `json:"patient_id"`
	ProviderID  string        `json:"provider_id"`
	ServiceType string        `json:"service_type"`
	StartTime   time.Time     `json:"start_time"`
	EndTime     time.Time     `json:"end_time"`
	Status      BookingStatus `json:"status"`
	Notes       string        `json:"notes,omitempty"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
}

// CreateBookingRequest represents the input for creating a booking
type CreateBookingRequest struct {
	PatientID   string    `json:"patient_id" binding:"required"`
	ProviderID  string    `json:"provider_id" binding:"required"`
	ServiceType string    `json:"service_type" binding:"required"`
	StartTime   time.Time `json:"start_time" binding:"required"`
	EndTime     time.Time `json:"end_time" binding:"required"`
	Notes       string    `json:"notes"`
}

// BookingRepository defines the interface for booking persistence
type BookingRepository interface {
	Create(ctx context.Context, booking *Booking) error
	GetByID(ctx context.Context, id string) (*Booking, error)
	GetByPatientID(ctx context.Context, patientID string, limit, offset int) ([]*Booking, error)
	GetByProviderID(ctx context.Context, providerID string, limit, offset int) ([]*Booking, error)
	UpdateStatus(ctx context.Context, id string, status BookingStatus) error
	CheckAvailability(ctx context.Context, providerID string, startTime, endTime time.Time) (bool, error)
}

// BookingUsecase defines the business logic interface
type BookingUsecase interface {
	CreateBooking(ctx context.Context, req *CreateBookingRequest) (*Booking, error)
	GetBooking(ctx context.Context, id string) (*Booking, error)
	GetPatientBookings(ctx context.Context, patientID string, limit, offset int) ([]*Booking, error)
	GetProviderBookings(ctx context.Context, providerID string, limit, offset int) ([]*Booking, error)
	CancelBooking(ctx context.Context, id, userID string) error
	ConfirmBooking(ctx context.Context, id string) error
}
