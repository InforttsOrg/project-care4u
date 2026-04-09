package config

import "os"

type Config struct {
	Port        string
	DatabaseURL string
	NatsURL     string // Added NATS URL
	JwtSecret   string
}

func Load() *Config {
	return &Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:pass@db:5432/auth_db?sslmode=disable"),
		NatsURL:     getEnv("NATS_URL", "nats://nats:4222"),
		JwtSecret:   getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
