package usecase

import (
	"context"
	"fmt"

	"github.com/care4u/services/location/internal/domain"
)

type locationUsecase struct {
	repo domain.LocationRepository
}

func NewLocationUsecase(repo domain.LocationRepository) domain.LocationUsecase {
	return &locationUsecase{repo: repo}
}

func (u *locationUsecase) UpdateLocation(ctx context.Context, req *domain.UpdateLocationRequest) error {
	// Validate coordinates
	if req.Latitude < -90 || req.Latitude > 90 {
		return fmt.Errorf("invalid latitude: must be between -90 and 90")
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		return fmt.Errorf("invalid longitude: must be between -180 and 180")
	}

	return u.repo.UpdateLocation(ctx, req)
}

func (u *locationUsecase) GetNearbyProviders(ctx context.Context, req *domain.NearbyRequest) ([]*domain.Provider, error) {
	// Validate and set defaults
	if req.RadiusKM <= 0 {
		req.RadiusKM = 10 // 10km default
	}
	if req.RadiusKM > 100 {
		req.RadiusKM = 100 // Max 100km
	}
	if req.Limit <= 0 {
		req.Limit = 20
	}

	return u.repo.GetNearbyProviders(ctx, req)
}
