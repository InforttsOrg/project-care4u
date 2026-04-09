package http

import (
	"net/http"

	"github.com/care4u/services/payment/internal/domain"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	usecase domain.PaymentUsecase
}

func NewPaymentHandler(usecase domain.PaymentUsecase) *PaymentHandler {
	return &PaymentHandler{usecase: usecase}
}

func (h *PaymentHandler) RegisterRoutes(r *gin.Engine) {
	payments := r.Group("/payments")
	{
		payments.POST("/initiate", h.InitiatePayment)
		payments.POST("/webhook", h.HandleWebhook)
		payments.GET("/:id/status", h.GetPaymentStatus)
	}
}

func (h *PaymentHandler) InitiatePayment(c *gin.Context) {
	var req domain.InitiatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.usecase.InitiatePayment(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	// In production, verify X-Razorpay-Signature header
	
	var payload domain.WebhookPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.HandleWebhook(c.Request.Context(), &payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

func (h *PaymentHandler) GetPaymentStatus(c *gin.Context) {
	id := c.Param("id")
	
	payment, err := h.usecase.GetPaymentStatus(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, payment)
}
