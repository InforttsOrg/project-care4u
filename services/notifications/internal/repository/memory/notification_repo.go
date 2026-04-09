package memory

import (
	"context"
	"fmt"
	"sync"

	"github.com/care4u/services/notifications/internal/domain"
)

// In-memory repository for MVP - replace with PostgreSQL or DynamoDB in production
type notificationRepository struct {
	mu      sync.RWMutex
	data    map[string]*domain.Notification
	byUser  map[string][]string // userID -> []notificationIDs
}

func NewNotificationRepository() domain.NotificationRepository {
	return &notificationRepository{
		data:   make(map[string]*domain.Notification),
		byUser: make(map[string][]string),
	}
}

func (r *notificationRepository) Create(ctx context.Context, notification *domain.Notification) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[notification.ID] = notification
	r.byUser[notification.UserID] = append(r.byUser[notification.UserID], notification.ID)
	return nil
}

func (r *notificationRepository) GetByID(ctx context.Context, id string) (*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if n, exists := r.data[id]; exists {
		return n, nil
	}
	return nil, nil
}

func (r *notificationRepository) GetByUserID(ctx context.Context, userID string, limit, offset int) ([]*domain.Notification, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := r.byUser[userID]
	if len(ids) == 0 {
		return []*domain.Notification{}, nil
	}

	// Apply offset and limit
	start := offset
	if start > len(ids) {
		return []*domain.Notification{}, nil
	}
	end := start + limit
	if end > len(ids) {
		end = len(ids)
	}

	// Return in reverse order (newest first)
	result := make([]*domain.Notification, 0, end-start)
	for i := len(ids) - 1 - start; i >= len(ids)-end && i >= 0; i-- {
		if n, exists := r.data[ids[i]]; exists {
			result = append(result, n)
		}
	}

	return result, nil
}

func (r *notificationRepository) UpdateStatus(ctx context.Context, id string, status domain.NotificationStatus, failReason string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if n, exists := r.data[id]; exists {
		n.Status = status
		n.FailReason = failReason
		return nil
	}
	return fmt.Errorf("notification not found")
}
