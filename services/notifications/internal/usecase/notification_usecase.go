package usecase

import (
	"context"
	"log"
	"time"

	"github.com/care4u/services/notifications/internal/domain"
	"github.com/google/uuid"
)

type notificationUsecase struct {
	repo        domain.NotificationRepository
	pushService domain.PushService
}

func NewNotificationUsecase(repo domain.NotificationRepository, pushService domain.PushService) domain.NotificationUsecase {
	return &notificationUsecase{
		repo:        repo,
		pushService: pushService,
	}
}

func (u *notificationUsecase) SendNotification(ctx context.Context, req *domain.SendNotificationRequest) (*domain.Notification, error) {
	now := time.Now()
	notification := &domain.Notification{
		ID:        uuid.New().String(),
		UserID:    req.UserID,
		Type:      req.Type,
		Title:     req.Title,
		Body:      req.Body,
		Data:      req.Data,
		Status:    domain.StatusQueued,
		CreatedAt: now,
	}

	if err := u.repo.Create(ctx, notification); err != nil {
		return nil, err
	}

	// Send notification based on type
	go u.processNotification(notification)

	return notification, nil
}

func (u *notificationUsecase) processNotification(n *domain.Notification) {
	ctx := context.Background()
	var err error

	switch n.Type {
	case domain.TypePush:
		if u.pushService != nil {
			err = u.pushService.Send(ctx, n.UserID, n.Title, n.Body, n.Data)
		} else {
			log.Printf("MOCK: Push notification to %s: %s - %s", n.UserID, n.Title, n.Body)
		}
	case domain.TypeSMS:
		log.Printf("MOCK: SMS to %s: %s", n.UserID, n.Body)
	case domain.TypeEmail:
		log.Printf("MOCK: Email to %s: %s - %s", n.UserID, n.Title, n.Body)
	case domain.TypeInApp:
		log.Printf("MOCK: In-app notification for %s: %s", n.UserID, n.Title)
	}

	now := time.Now()
	if err != nil {
		u.repo.UpdateStatus(ctx, n.ID, domain.StatusFailed, err.Error())
	} else {
		n.SentAt = &now
		u.repo.UpdateStatus(ctx, n.ID, domain.StatusSent, "")
	}
}

func (u *notificationUsecase) GetUserNotifications(ctx context.Context, userID string, limit, offset int) ([]*domain.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	return u.repo.GetByUserID(ctx, userID, limit, offset)
}
