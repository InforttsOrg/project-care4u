package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/care4u/services/auth/internal/config"
	"github.com/care4u/services/auth/internal/domain"
	"github.com/care4u/services/auth/internal/event"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authUsecase struct {
	repo     domain.UserRepository
	producer event.EventProducer
	config   *config.Config
}

func NewAuthUsecase(repo domain.UserRepository, producer event.EventProducer, cfg *config.Config) domain.AuthUsecase {
	return &authUsecase{
		repo:     repo,
		producer: producer,
		config:   cfg,
	}
}

func (u *authUsecase) SendOTP(ctx context.Context, phone string) error {
	// Mock OTP sending
	fmt.Printf("MOCK OTP for %s: 1234\n", phone)
	return nil
}

func (u *authUsecase) VerifyOTP(ctx context.Context, phone, otp string) (string, error) {
	// Mock verification
	if otp != "1234" {
		return "", fmt.Errorf("invalid OTP")
	}

	// Check if user exists, else create
	user, err := u.repo.FindByPhone(ctx, phone)
	if err != nil {
		return "", err
	}

	if user == nil {
		user = &domain.User{
			ID:        uuid.New().String(),
			Phone:     phone,
			Role:      "patient", // Default role
			CreatedAt: time.Now(),
		}
		if err := u.repo.Create(ctx, user); err != nil {
			return "", err
		}
	}

	// Publish Event
	if err := u.producer.PublishUserVerified(user); err != nil {
		fmt.Printf("Warning: Failed to publish event: %v\n", err)
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(u.config.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
