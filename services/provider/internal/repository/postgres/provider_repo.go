package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/care4u/services/provider/internal/domain"
	"github.com/lib/pq"
)

type providerRepository struct {
	db *sql.DB
}

func NewProviderRepository(db *sql.DB) domain.ProviderRepository {
	return &providerRepository{db: db}
}

func (r *providerRepository) Create(ctx context.Context, provider *domain.Provider) error {
	query := `
		INSERT INTO providers (id, user_id, name, specialty, qualifications, experience_years, bio, consultation_fee, rating, review_count, status, is_available, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		provider.ID,
		provider.UserID,
		provider.Name,
		provider.Specialty,
		pq.Array(provider.Qualifications),
		provider.Experience,
		provider.Bio,
		provider.ConsultationFee,
		provider.Rating,
		provider.ReviewCount,
		provider.Status,
		provider.IsAvailable,
		provider.CreatedAt,
		provider.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create provider: %w", err)
	}
	return nil
}

func (r *providerRepository) GetByID(ctx context.Context, id string) (*domain.Provider, error) {
	return r.scanProvider(ctx, "SELECT id, user_id, name, specialty, qualifications, experience_years, bio, consultation_fee, rating, review_count, status, is_available, created_at, updated_at FROM providers WHERE id = $1", id)
}

func (r *providerRepository) GetByUserID(ctx context.Context, userID string) (*domain.Provider, error) {
	return r.scanProvider(ctx, "SELECT id, user_id, name, specialty, qualifications, experience_years, bio, consultation_fee, rating, review_count, status, is_available, created_at, updated_at FROM providers WHERE user_id = $1", userID)
}

func (r *providerRepository) scanProvider(ctx context.Context, query, arg string) (*domain.Provider, error) {
	row := r.db.QueryRowContext(ctx, query, arg)

	var p domain.Provider
	var qualifications []string

	err := row.Scan(
		&p.ID,
		&p.UserID,
		&p.Name,
		&p.Specialty,
		pq.Array(&qualifications),
		&p.Experience,
		&p.Bio,
		&p.ConsultationFee,
		&p.Rating,
		&p.ReviewCount,
		&p.Status,
		&p.IsAvailable,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan provider: %w", err)
	}
	p.Qualifications = qualifications
	return &p, nil
}

func (r *providerRepository) Update(ctx context.Context, id string, req *domain.UpdateProviderRequest) error {
	// Dynamic update query - simplified for MVP
	query := `UPDATE providers SET updated_at = NOW()`
	args := []interface{}{}
	argNum := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argNum)
		args = append(args, *req.Name)
		argNum++
	}
	if req.Specialty != nil {
		query += fmt.Sprintf(", specialty = $%d", argNum)
		args = append(args, *req.Specialty)
		argNum++
	}
	if req.Experience != nil {
		query += fmt.Sprintf(", experience_years = $%d", argNum)
		args = append(args, *req.Experience)
		argNum++
	}
	if req.Bio != nil {
		query += fmt.Sprintf(", bio = $%d", argNum)
		args = append(args, *req.Bio)
		argNum++
	}
	if req.ConsultationFee != nil {
		query += fmt.Sprintf(", consultation_fee = $%d", argNum)
		args = append(args, *req.ConsultationFee)
		argNum++
	}
	if req.IsAvailable != nil {
		query += fmt.Sprintf(", is_available = $%d", argNum)
		args = append(args, *req.IsAvailable)
		argNum++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argNum)
	args = append(args, id)

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update provider: %w", err)
	}
	return nil
}

func (r *providerRepository) Verify(ctx context.Context, id string) error {
	// Verification is a dedicated write: it flips status to 'verified' AND makes the
	// profile visible. Search() only returns status='verified' providers, so an
	// availability-only update would leave verified providers permanently hidden.
	query := `UPDATE providers SET status = $1, is_available = true, updated_at = NOW() WHERE id = $2`
	result, err := r.db.ExecContext(ctx, query, domain.ProviderStatusVerified, id)
	if err != nil {
		return fmt.Errorf("failed to verify provider: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("provider not found")
	}
	return nil
}

func (r *providerRepository) Search(ctx context.Context, specialty string, limit, offset int) ([]*domain.Provider, error) {
	query := `
		SELECT id, user_id, name, specialty, qualifications, experience_years, bio, consultation_fee, rating, review_count, status, is_available, created_at, updated_at 
		FROM providers 
		WHERE status = 'verified' AND is_available = true
	`
	args := []interface{}{}
	argNum := 1

	if specialty != "" {
		query += fmt.Sprintf(" AND specialty ILIKE $%d", argNum)
		args = append(args, "%"+specialty+"%")
		argNum++
	}

	query += fmt.Sprintf(" ORDER BY rating DESC LIMIT $%d OFFSET $%d", argNum, argNum+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search providers: %w", err)
	}
	defer rows.Close()

	var providers []*domain.Provider
	for rows.Next() {
		var p domain.Provider
		var qualifications []string
		err := rows.Scan(
			&p.ID,
			&p.UserID,
			&p.Name,
			&p.Specialty,
			pq.Array(&qualifications),
			&p.Experience,
			&p.Bio,
			&p.ConsultationFee,
			&p.Rating,
			&p.ReviewCount,
			&p.Status,
			&p.IsAvailable,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider: %w", err)
		}
		p.Qualifications = qualifications
		providers = append(providers, &p)
	}
	return providers, rows.Err()
}

func (r *providerRepository) UpdateRating(ctx context.Context, id string, rating float64, reviewCount int) error {
	query := `UPDATE providers SET rating = $1, review_count = $2, updated_at = NOW() WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, rating, reviewCount, id)
	if err != nil {
		return fmt.Errorf("failed to update rating: %w", err)
	}
	return nil
}
