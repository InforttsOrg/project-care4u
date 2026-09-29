package http

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/care4u/services/payment/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

const testWebhookSecret = "whsec_test_secret"

type stubPaymentUsecase struct {
	handled   []*domain.WebhookPayload
	initiated int
}

func (s *stubPaymentUsecase) HandleWebhook(_ context.Context, payload *domain.WebhookPayload) error {
	s.handled = append(s.handled, payload)
	return nil
}

func (s *stubPaymentUsecase) InitiatePayment(_ context.Context, _ *domain.InitiatePaymentRequest) (*domain.InitiatePaymentResponse, error) {
	s.initiated++
	return &domain.InitiatePaymentResponse{PaymentID: "pay_1", OrderID: "order_1"}, nil
}

func (s *stubPaymentUsecase) GetPaymentStatus(_ context.Context, _ string) (*domain.Payment, error) {
	return &domain.Payment{ID: "pay_1"}, nil
}

func sign(body, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(body))
	return hex.EncodeToString(mac.Sum(nil))
}

// postWebhook drives POST /payments/webhook and returns the recorder plus the
// usecase stub so a test can assert the payload never reached the usecase.
func postWebhook(t *testing.T, secret, body, signature string) (*httptest.ResponseRecorder, *stubPaymentUsecase) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	usecase := &stubPaymentUsecase{}
	r := gin.New()
	NewPaymentHandler(usecase, secret).RegisterRoutes(r)

	req := httptest.NewRequest("POST", "/payments/webhook", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if signature != "" {
		req.Header.Set("X-Razorpay-Signature", signature)
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w, usecase
}

const capturedBody = `{"entity":"event","event":"payment.captured","payload":{"payment":{"entity":{"id":"pay_1","order_id":"order_1"}}}}`

func TestWebhookRejectsBadSignature(t *testing.T) {
	cases := map[string]struct{ secret, body, signature string }{
		"wrong signature": {testWebhookSecret, capturedBody, sign(capturedBody, "other_secret")},
		"missing header":  {testWebhookSecret, capturedBody, ""},
		"blank signature": {testWebhookSecret, capturedBody, "  "},
		"not hex":         {testWebhookSecret, capturedBody, "zzzz"},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			w, usecase := postWebhook(t, tc.secret, tc.body, tc.signature)
			if w.Code != 401 {
				t.Fatalf("status = %d, want 401 (body %s)", w.Code, w.Body.String())
			}
			if len(usecase.handled) != 0 {
				t.Fatalf("rejected webhook must not reach the usecase, got %d", len(usecase.handled))
			}
			if strings.Contains(w.Body.String(), testWebhookSecret) {
				t.Fatal("error response must not echo the webhook secret")
			}
		})
	}
}

func TestWebhookRejectsTamperedBody(t *testing.T) {
	// A valid signature over a different body must not be accepted.
	w, usecase := postWebhook(t, testWebhookSecret, capturedBody, sign(`{"event":"payment.captured"}`, testWebhookSecret))
	if w.Code != 401 {
		t.Fatalf("status = %d, want 401", w.Code)
	}
	if len(usecase.handled) != 0 {
		t.Fatalf("tampered webhook must not reach the usecase, got %d", len(usecase.handled))
	}
}

func TestWebhookAcceptsValidSignature(t *testing.T) {
	w, usecase := postWebhook(t, testWebhookSecret, capturedBody, sign(capturedBody, testWebhookSecret))
	if w.Code != 200 {
		t.Fatalf("status = %d, want 200 (body %s)", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "received") {
		t.Fatalf("unexpected body %s", w.Body.String())
	}
	if len(usecase.handled) != 1 {
		t.Fatalf("usecase must handle the webhook exactly once, got %d", len(usecase.handled))
	}
	if usecase.handled[0].Event != "payment.captured" {
		t.Fatalf("parsed event = %q, want payment.captured", usecase.handled[0].Event)
	}
}

// With no secret configured (local/mock flow) the endpoint must behave exactly
// as it did before signature verification existed.
func TestWebhookWithoutSecretStaysOpen(t *testing.T) {
	w, usecase := postWebhook(t, "", capturedBody, "")
	if w.Code != 200 {
		t.Fatalf("status = %d, want 200 (body %s)", w.Code, w.Body.String())
	}
	if len(usecase.handled) != 1 {
		t.Fatalf("mock-mode webhook must reach the usecase, got %d", len(usecase.handled))
	}
}

func TestWebhookRejectsMalformedBody(t *testing.T) {
	body := `{"event":`
	w, usecase := postWebhook(t, testWebhookSecret, body, sign(body, testWebhookSecret))
	if w.Code != 400 {
		t.Fatalf("status = %d, want 400", w.Code)
	}
	if len(usecase.handled) != 0 {
		t.Fatal("malformed webhook must not reach the usecase")
	}
}

func TestWebhookRejectsOversizedBody(t *testing.T) {
	body := `{"event":"payment.captured","padding":"` + strings.Repeat("a", maxWebhookBodyBytes+16) + `"}`
	w, _ := postWebhook(t, testWebhookSecret, body, sign(body, testWebhookSecret))
	if w.Code != 400 {
		t.Fatalf("status = %d, want 400 for a body over the cap", w.Code)
	}
}
