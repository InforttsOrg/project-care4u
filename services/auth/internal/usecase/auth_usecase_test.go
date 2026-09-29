package usecase

import (
	"context"
	"strings"
	"testing"

	"github.com/care4u/services/auth/internal/config"
	"github.com/care4u/services/auth/internal/domain"
)

// stubUserRepo is an in-memory UserRepository for the usecase tests.
type stubUserRepo struct {
	users   map[string]*domain.User
	created int
}

func newStubUserRepo() *stubUserRepo {
	return &stubUserRepo{users: map[string]*domain.User{}}
}

func (s *stubUserRepo) Create(_ context.Context, user *domain.User) error {
	if s.users == nil {
		s.users = map[string]*domain.User{}
	}
	s.users[user.Phone] = user
	s.created++
	return nil
}

func (s *stubUserRepo) FindByPhone(_ context.Context, phone string) (*domain.User, error) {
	return s.users[phone], nil
}

// recordingProducer stands in for the NATS producer.
type recordingProducer struct {
	published []*domain.User
}

func (p *recordingProducer) PublishUserVerified(user *domain.User) error {
	p.published = append(p.published, user)
	return nil
}

func testConfig() *config.Config {
	return &config.Config{
		Port:        "8080",
		DatabaseURL: "postgres://user:pass@localhost:5432/auth_db?sslmode=disable",
		NatsURL:     "nats://localhost:4222",
		JwtSecret:   "unit-test-secret",
	}
}

func TestVerifyOTPRejectsWrongCode(t *testing.T) {
	repo := newStubUserRepo()
	u := NewAuthUsecase(repo, &recordingProducer{}, testConfig())

	if _, err := u.VerifyOTP(context.Background(), "+919999999999", "0000"); err == nil {
		t.Fatal("expected an error for a wrong OTP")
	}
	if repo.created != 0 {
		t.Fatalf("a rejected OTP must not create a user, created=%d", repo.created)
	}
}

func TestVerifyOTPProvisionsUserAndReturnsToken(t *testing.T) {
	repo := newStubUserRepo()
	producer := &recordingProducer{}
	u := NewAuthUsecase(repo, producer, testConfig())

	token, err := u.VerifyOTP(context.Background(), "+919999999999", "1234")
	if err != nil {
		t.Fatalf("VerifyOTP: %v", err)
	}
	if strings.Count(token, ".") != 2 {
		t.Fatalf("expected a JWT (header.payload.signature), got %q", token)
	}
	if repo.created != 1 {
		t.Fatalf("first-time login must provision exactly one user, created=%d", repo.created)
	}
	user := repo.users["+919999999999"]
	if user == nil {
		t.Fatal("user was not stored")
	}
	if user.Role != "patient" {
		t.Fatalf("default role = %q, want patient", user.Role)
	}
	if len(producer.published) != 1 {
		t.Fatalf("user.verified must be published once, got %d", len(producer.published))
	}

	// Second login reuses the stored user and still issues a token.
	if _, err := u.VerifyOTP(context.Background(), "+919999999999", "1234"); err != nil {
		t.Fatalf("second VerifyOTP: %v", err)
	}
	if repo.created != 1 {
		t.Fatalf("returning user must not be re-created, created=%d", repo.created)
	}
}

// cmd/api/main.go keeps serving when NATS is unavailable: the producer is nil and
// the OTP login path used to dereference it (nil interface call -> panic -> 500).
func TestVerifyOTPSurvivesUnavailableProducer(t *testing.T) {
	repo := newStubUserRepo()
	u := NewAuthUsecase(repo, nil, testConfig())

	token, err := u.VerifyOTP(context.Background(), "+919999999999", "1234")
	if err != nil {
		t.Fatalf("VerifyOTP with no event producer: %v", err)
	}
	if token == "" {
		t.Fatal("expected a token even when the producer is unavailable")
	}
	if repo.created != 1 {
		t.Fatalf("user must still be provisioned, created=%d", repo.created)
	}
}

func TestSendOTP(t *testing.T) {
	u := NewAuthUsecase(newStubUserRepo(), &recordingProducer{}, testConfig())
	if err := u.SendOTP(context.Background(), "+919999999999"); err != nil {
		t.Fatalf("SendOTP: %v", err)
	}
}
