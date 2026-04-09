package main

import (
	"log"

	"github.com/care4u/services/location/internal/config"
	redisrepo "github.com/care4u/services/location/internal/repository/redis"
	transport "github.com/care4u/services/location/internal/transport/http"
	"github.com/care4u/services/location/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

func main() {
	cfg := config.Load()

	// Redis Connection
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})

	// Dependency Injection
	locationRepo := redisrepo.NewLocationRepository(rdb)
	locationUsecase := usecase.NewLocationUsecase(locationRepo)
	locationHandler := transport.NewLocationHandler(locationUsecase)

	// Router Setup
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "up", "service": "location"})
	})

	locationHandler.RegisterRoutes(r)

	// Start Server
	log.Printf("Location service starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
