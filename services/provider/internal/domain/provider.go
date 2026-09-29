package domain

import (
	"context"
	"time"
)

// ProviderStatus represents provider verification status
type ProviderStatus string

const (
	ProviderStatusPending   ProviderStatus = "pending"
	ProviderStatusVerified  ProviderStatus = "verified"
	ProviderStatusSuspended ProviderStatus = "suspended"
)

// Provider represents a healthcare provider (doctor, clinic, etc.)
type Provider struct {
	ID              string         `json:"id"`
	UserID          string         `json:"user_id"` // Links to auth service user
	Name            string         `json:"name"`
	Specialty       string         `json:"specialty"`
	Qualifications  []string       `json:"qualifications"`
	Experience      int            `json:"experience_years"`
	Bio             string         `json:"bio,omitempty"`
	ConsultationFee int64          `json:"consultation_fee"` // In paise
	Rating          float64        `json:"rating"`
	ReviewCount     int            `json:"review_count"`
	Status          ProviderStatus `json:"status"`
	IsAvailable     bool           `json:"is_available"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

// TimeSlot represents an available time slot
type TimeSlot struct {
	ID           string `json:"id"`
	ProviderID   string `json:"provider_id"`
	DayOfWeek    int    `json:"day_of_week"` // 0=Sunday, 6=Saturday
	StartTime    string `json:"start_time"`  // "09:00"
	EndTime      string `json:"end_time"`    // "17:00"
	SlotDuration int    `json:"slot_duration_minutes"`
}

// CreateProviderRequest represents the input for creating a provider profile
type CreateProviderRequest struct {
	UserID          string   `json:"user_id" binding:"required"`
	Name            string   `json:"name" binding:"required"`
	Specialty       string   `json:"specialty" binding:"required"`
	Qualifications  []string `json:"qualifications"`
	Experience      int      `json:"experience_years"`
	Bio             string   `json:"bio"`
	ConsultationFee int64    `json:"consultation_fee"`
}

// UpdateProviderRequest represents the input for updating a provider profile
type UpdateProviderRequest struct {
	Name            *string  `json:"name"`
	Specialty       *string  `json:"specialty"`
	Qualifications  []string `json:"qualifications"`
	Experience      *int     `json:"experience_years"`
	Bio             *string  `json:"bio"`
	ConsultationFee *int64   `json:"consultation_fee"`
	IsAvailable     *bool    `json:"is_available"`
}

// ProviderRepository defines persistence operations
type ProviderRepository interface {
	Create(ctx context.Context, provider *Provider) error
	GetByID(ctx context.Context, id string) (*Provider, error)
	GetByUserID(ctx context.Context, userID string) (*Provider, error)
	Update(ctx context.Context, id string, req *UpdateProviderRequest) error
	Search(ctx context.Context, specialty string, limit, offset int) ([]*Provider, error)
	UpdateRating(ctx context.Context, id string, rating float64, reviewCount int) error
}

// ProviderUsecase defines business logic
type ProviderUsecase interface {
	CreateProvider(ctx context.Context, req *CreateProviderRequest) (*Provider, error)
	GetProvider(ctx context.Context, id string) (*Provider, error)
	GetProviderByUserID(ctx context.Context, userID string) (*Provider, error)
	UpdateProvider(ctx context.Context, id string, req *UpdateProviderRequest) error
	SearchProviders(ctx context.Context, specialty string, limit, offset int) ([]*Provider, error)
	VerifyProvider(ctx context.Context, id string) error
}
