package domain

import (
	"context"
	"time"
)

type User struct {
	ID        string    `json:"id"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"` // "patient" or "provider"
	CreatedAt time.Time `json:"created_at"`
}

// Repository Interface
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByPhone(ctx context.Context, phone string) (*User, error)
}

// Usecase Interface
type AuthUsecase interface {
	SendOTP(ctx context.Context, phone string) error
	VerifyOTP(ctx context.Context, phone, otp string) (string, error) // Returns JWT
}
