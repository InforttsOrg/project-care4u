package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/care4u/services/auth/internal/domain"
	_ "github.com/lib/pq"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (id, phone, role, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Phone, user.Role, user.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (*domain.User, error) {
	query := `SELECT id, phone, role, created_at FROM users WHERE phone = $1`
	row := r.db.QueryRowContext(ctx, query, phone)

	var user domain.User
	err := row.Scan(&user.ID, &user.Phone, &user.Role, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Not found
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}
	return &user, nil
}
