package config

import "os"

type Config struct {
	Port     string
	RedisURL string
}

func Load() *Config {
	return &Config{
		Port:     getEnv("PORT", "8083"),
		RedisURL: getEnv("REDIS_URL", "redis:6379"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
