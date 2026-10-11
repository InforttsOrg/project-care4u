package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/care4u/services/provider/internal/domain"
	"github.com/google/uuid"
)

type providerUsecase struct {
	repo domain.ProviderRepository
}

func NewProviderUsecase(repo domain.ProviderRepository) domain.ProviderUsecase {
	return &providerUsecase{repo: repo}
}

func (u *providerUsecase) CreateProvider(ctx context.Context, req *domain.CreateProviderRequest) (*domain.Provider, error) {
	// Check if provider profile already exists for this user
	existing, err := u.repo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("provider profile already exists for this user")
	}

	now := time.Now()
	provider := &domain.Provider{
		ID:              uuid.New().String(),
		UserID:          req.UserID,
		Name:            req.Name,
		Specialty:       req.Specialty,
		Qualifications:  req.Qualifications,
		Experience:      req.Experience,
		Bio:             req.Bio,
		ConsultationFee: req.ConsultationFee,
		Rating:          0,
		ReviewCount:     0,
		Status:          domain.ProviderStatusPending,
		IsAvailable:     false, // Initially unavailable until verified
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := u.repo.Create(ctx, provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (u *providerUsecase) GetProvider(ctx context.Context, id string) (*domain.Provider, error) {
	provider, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, fmt.Errorf("provider not found")
	}
	return provider, nil
}

func (u *providerUsecase) GetProviderByUserID(ctx context.Context, userID string) (*domain.Provider, error) {
	provider, err := u.repo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, fmt.Errorf("provider not found")
	}
	return provider, nil
}

func (u *providerUsecase) UpdateProvider(ctx context.Context, id string, req *domain.UpdateProviderRequest) error {
	provider, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("provider not found")
	}

	return u.repo.Update(ctx, id, req)
}

func (u *providerUsecase) SearchProviders(ctx context.Context, specialty string, limit, offset int) ([]*domain.Provider, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.Search(ctx, specialty, limit, offset)
}

func (u *providerUsecase) VerifyProvider(ctx context.Context, id string) error {
	provider, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if provider == nil {
		return fmt.Errorf("provider not found")
	}

	// Flip the profile to verified and make it visible. This goes through a
	// dedicated repo write, not the generic Update (which has no status field on
	// purpose: verification must not be reachable via the public PUT /providers/:id payload).
	return u.repo.Verify(ctx, id)
}
