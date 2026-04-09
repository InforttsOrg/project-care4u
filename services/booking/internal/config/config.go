package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	NatsURL     string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@db:5432/booking_db?sslmode=disable"),
		NatsURL:     getEnv("NATS_URL", "nats://nats:4222"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
