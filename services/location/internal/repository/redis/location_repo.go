package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/care4u/services/location/internal/domain"
	"github.com/go-redis/redis/v8"
)

const (
	providerLocationKey = "provider:location:"
	geoKey              = "providers:geo"
)

type locationRepository struct {
	client *redis.Client
}

func NewLocationRepository(client *redis.Client) domain.LocationRepository {
	return &locationRepository{client: client}
}

func (r *locationRepository) UpdateLocation(ctx context.Context, req *domain.UpdateLocationRequest) error {
	// Store in Redis GEO
	err := r.client.GeoAdd(ctx, geoKey, &redis.GeoLocation{
		Name:      req.ProviderID,
		Longitude: req.Longitude,
		Latitude:  req.Latitude,
	}).Err()
	if err != nil {
		return fmt.Errorf("failed to update geo location: %w", err)
	}

	// Store additional provider data
	provider := domain.Provider{
		ID:          req.ProviderID,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Address:     req.Address,
		LastUpdated: time.Now(),
	}

	data, err := json.Marshal(provider)
	if err != nil {
		return fmt.Errorf("failed to marshal provider: %w", err)
	}

	return r.client.Set(ctx, providerLocationKey+req.ProviderID, data, 24*time.Hour).Err()
}

func (r *locationRepository) GetNearbyProviders(ctx context.Context, req *domain.NearbyRequest) ([]*domain.Provider, error) {
	radius := req.RadiusKM
	if radius <= 0 {
		radius = 10 // Default 10km
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	// Query GEO
	locations, err := r.client.GeoRadius(ctx, geoKey, req.Longitude, req.Latitude, &redis.GeoRadiusQuery{
		Radius:    radius,
		Unit:      "km",
		WithCoord: true,
		WithDist:  true,
		Count:     limit,
		Sort:      "ASC",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to query nearby providers: %w", err)
	}

	providers := make([]*domain.Provider, 0, len(locations))
	for _, loc := range locations {
		provider := &domain.Provider{
			ID:        loc.Name,
			Latitude:  loc.Latitude,
			Longitude: loc.Longitude,
			Distance:  loc.Dist,
		}

		// Fetch additional data
		data, err := r.client.Get(ctx, providerLocationKey+loc.Name).Bytes()
		if err == nil {
			var stored domain.Provider
			if json.Unmarshal(data, &stored) == nil {
				provider.Name = stored.Name
				provider.Specialty = stored.Specialty
				provider.Address = stored.Address
				provider.Rating = stored.Rating
				provider.LastUpdated = stored.LastUpdated
			}
		}

		providers = append(providers, provider)
	}

	return providers, nil
}

func (r *locationRepository) GetProviderLocation(ctx context.Context, providerID string) (*domain.Provider, error) {
	data, err := r.client.Get(ctx, providerLocationKey+providerID).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get provider location: %w", err)
	}

	var provider domain.Provider
	if err := json.Unmarshal(data, &provider); err != nil {
		return nil, fmt.Errorf("failed to unmarshal provider: %w", err)
	}

	return &provider, nil
}
