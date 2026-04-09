package domain

import (
	"context"
	"time"
)

// Provider represents a healthcare provider with location
type Provider struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Specialty   string    `json:"specialty"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	Address     string    `json:"address,omitempty"`
	Rating      float64   `json:"rating"`
	Distance    float64   `json:"distance,omitempty"` // Calculated field
	LastUpdated time.Time `json:"last_updated"`
}

// UpdateLocationRequest represents a location update
type UpdateLocationRequest struct {
	ProviderID string  `json:"provider_id" binding:"required"`
	Latitude   float64 `json:"latitude" binding:"required"`
	Longitude  float64 `json:"longitude" binding:"required"`
	Address    string  `json:"address"`
}

// NearbyRequest represents a search for nearby providers
type NearbyRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
	RadiusKM  float64 `json:"radius_km"` // Default 10km
	Specialty string  `json:"specialty"`
	Limit     int     `json:"limit"`
}

// LocationRepository defines location persistence
type LocationRepository interface {
	UpdateLocation(ctx context.Context, req *UpdateLocationRequest) error
	GetNearbyProviders(ctx context.Context, req *NearbyRequest) ([]*Provider, error)
	GetProviderLocation(ctx context.Context, providerID string) (*Provider, error)
}

// LocationUsecase defines business logic
type LocationUsecase interface {
	UpdateLocation(ctx context.Context, req *UpdateLocationRequest) error
	GetNearbyProviders(ctx context.Context, req *NearbyRequest) ([]*Provider, error)
}
