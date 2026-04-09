package http

import (
	"net/http"
	"strconv"

	"github.com/care4u/services/notifications/internal/domain"
	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	usecase domain.NotificationUsecase
}

func NewNotificationHandler(usecase domain.NotificationUsecase) *NotificationHandler {
	return &NotificationHandler{usecase: usecase}
}

func (h *NotificationHandler) RegisterRoutes(r *gin.Engine) {
	notifications := r.Group("/notifications")
	{
		notifications.POST("/send", h.SendNotification)
		notifications.GET("/history", h.GetHistory)
	}
}

func (h *NotificationHandler) SendNotification(c *gin.Context) {
	var req domain.SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notification, err := h.usecase.SendNotification(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":         "Notification queued",
		"notification_id": notification.ID,
	})
}

func (h *NotificationHandler) GetHistory(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	limit := 20
	offset := 0
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}
	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil && parsed >= 0 {
			offset = parsed
		}
	}

	notifications, err := h.usecase.GetUserNotifications(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications, "count": len(notifications)})
}
