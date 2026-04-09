package config

import "os"

type Config struct {
	Port           string
	NatsURL        string
	FCMCredentials string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "8084"),
		NatsURL:        getEnv("NATS_URL", "nats://nats:4222"),
		FCMCredentials: getEnv("FCM_CREDENTIALS", ""),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
