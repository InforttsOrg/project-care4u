package config

import "os"

type Config struct {
	Port              string
	DatabaseURL       string
	RazorpayKeyID     string
	RazorpayKeySecret string
}

func Load() *Config {
	return &Config{
		Port:              getEnv("PORT", "8082"),
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://user:pass@db:5432/payment_db?sslmode=disable"),
		RazorpayKeyID:     getEnv("RAZORPAY_KEY_ID", "rzp_test_xxx"),
		RazorpayKeySecret: getEnv("RAZORPAY_KEY_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
