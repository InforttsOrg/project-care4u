package http

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/care4u/services/payment/internal/domain"
	"github.com/gin-gonic/gin"
)

// maxWebhookBodyBytes caps the inbound webhook body. /payments/webhook is
// unauthenticated by design (it is called by the gateway, not a browser), so the
// body is read through a limit instead of into memory unbounded.
const maxWebhookBodyBytes = 1 << 20 // 1 MiB

type PaymentHandler struct {
	usecase domain.PaymentUsecase
	// webhookSecret, when set, is required on every inbound webhook. Empty
	// preserves the local/mock flow where the gateway secret is not configured.
	webhookSecret string
}

func NewPaymentHandler(usecase domain.PaymentUsecase, webhookSecret string) *PaymentHandler {
	return &PaymentHandler{usecase: usecase, webhookSecret: webhookSecret}
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
	// Razorpay signs the RAW request body, so the signature has to be checked
	// before the payload is parsed — and before any payment row is touched.
	body, err := io.ReadAll(http.MaxBytesReader(c.Writer, c.Request.Body, maxWebhookBodyBytes))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "could not read webhook body"})
		return
	}

	if !h.verifySignature(c.GetHeader("X-Razorpay-Signature"), body) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid webhook signature"})
		return
	}

	var payload domain.WebhookPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.usecase.HandleWebhook(c.Request.Context(), &payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "received"})
}

// verifySignature reports whether the request is authentic. With no webhook
// secret configured the endpoint keeps its mock behaviour (verified=false is
// never consulted for a nil secret) and warns once so the gap is visible in
// service logs rather than silent.
func (h *PaymentHandler) verifySignature(signature string, body []byte) bool {
	if h.webhookSecret == "" {
		log.Printf("Warning: RAZORPAY_WEBHOOK_SECRET is not set — /payments/webhook accepts unverified payloads")
		return true
	}
	if signature == "" {
		return false
	}
	mac := hmac.New(sha256.New, []byte(h.webhookSecret))
	mac.Write(body)
	expected := hex.EncodeToString(mac.Sum(nil))
	// Constant-time compare: a plain == would leak the signature prefix.
	return hmac.Equal([]byte(strings.ToLower(strings.TrimSpace(signature))), []byte(expected))
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
