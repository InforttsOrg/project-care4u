package usecase

import (
	"context"
	"fmt"
	"testing"

	"github.com/care4u/services/provider/internal/domain"
)

// stubProviderRepo is an in-memory ProviderRepository for usecase tests.
type stubProviderRepo struct {
	providers map[string]*domain.Provider
}

func newStubProviderRepo() *stubProviderRepo {
	return &stubProviderRepo{providers: map[string]*domain.Provider{}}
}

func (s *stubProviderRepo) Create(_ context.Context, p *domain.Provider) error {
	s.providers[p.ID] = p
	return nil
}

func (s *stubProviderRepo) GetByID(_ context.Context, id string) (*domain.Provider, error) {
	return s.providers[id], nil
}

func (s *stubProviderRepo) GetByUserID(_ context.Context, userID string) (*domain.Provider, error) {
	for _, p := range s.providers {
		if p.UserID == userID {
			return p, nil
		}
	}
	return nil, nil
}

func (s *stubProviderRepo) Update(_ context.Context, _ string, _ *domain.UpdateProviderRequest) error {
	return nil
}

func (s *stubProviderRepo) Verify(_ context.Context, id string) error {
	p, ok := s.providers[id]
	if !ok {
		return fmt.Errorf("provider not found")
	}
	p.Status = domain.ProviderStatusVerified
	p.IsAvailable = true
	return nil
}

func (s *stubProviderRepo) Search(_ context.Context, _ string, _, _ int) ([]*domain.Provider, error) {
	return nil, nil
}

func (s *stubProviderRepo) UpdateRating(_ context.Context, _ string, _ float64, _ int) error {
	return nil
}

// A verified provider must be status='verified': Search() filters on that column,
// so verifying used to only set is_available and left the profile invisible forever.
func TestVerifyProviderSetsVerifiedAndAvailable(t *testing.T) {
	repo := newStubProviderRepo()
	repo.providers["p1"] = &domain.Provider{
		ID:          "p1",
		Status:      domain.ProviderStatusPending,
		IsAvailable: false,
	}
	u := NewProviderUsecase(repo)

	if err := u.VerifyProvider(context.Background(), "p1"); err != nil {
		t.Fatalf("VerifyProvider: %v", err)
	}
	got := repo.providers["p1"]
	if got.Status != domain.ProviderStatusVerified {
		t.Fatalf("status = %q, want %q", got.Status, domain.ProviderStatusVerified)
	}
	if !got.IsAvailable {
		t.Fatal("is_available = false, want true after verification")
	}
}

func TestVerifyProviderRejectsUnknown(t *testing.T) {
	u := NewProviderUsecase(newStubProviderRepo())
	if err := u.VerifyProvider(context.Background(), "missing"); err == nil {
		t.Fatal("expected an error for an unknown provider")
	}
}

func TestCreateProviderStartsPendingAndUnavailable(t *testing.T) {
	repo := newStubProviderRepo()
	u := NewProviderUsecase(repo)

	p, err := u.CreateProvider(context.Background(), &domain.CreateProviderRequest{
		UserID:    "u1",
		Name:      "Dr. Test",
		Specialty: "cardiology",
	})
	if err != nil {
		t.Fatalf("CreateProvider: %v", err)
	}
	if p.Status != domain.ProviderStatusPending {
		t.Fatalf("status = %q, want pending", p.Status)
	}
	if p.IsAvailable {
		t.Fatal("a freshly created provider must not be available before verification")
	}
}
