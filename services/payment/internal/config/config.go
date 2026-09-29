package config

import "os"

type Config struct {
	Port              string
	DatabaseURL       string
	RazorpayKeyID     string
	RazorpayKeySecret string
	// RazorpayWebhookSecret signs inbound webhooks (X-Razorpay-Signature).
	// Empty keeps the local/mock flow: webhooks are then accepted unverified.
	RazorpayWebhookSecret string
}

func Load() *Config {
	return &Config{
		Port:                  getEnv("PORT", "8082"),
		DatabaseURL:           getEnv("DATABASE_URL", "postgres://user:pass@db:5432/payment_db?sslmode=disable"),
		RazorpayKeyID:         getEnv("RAZORPAY_KEY_ID", "rzp_test_xxx"),
		RazorpayKeySecret:     getEnv("RAZORPAY_KEY_SECRET", "secret"),
		RazorpayWebhookSecret: getEnv("RAZORPAY_WEBHOOK_SECRET", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
