package config

import (
	"context"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx   = context.Background()
	Redis *redis.Client
)

func InitRedis() error {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Println("DEBUG REDIS_URL:", redisURL)
		log.Fatal("REDIS_URL not set")
	}

	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Invalid REDIS_URL:", err)
	}

	Redis = redis.NewClient(opt)

	// Test connection
	_, err = Redis.Ping(Ctx).Result()
	if err != nil {
		return err
	}

	log.Println("✅ Redis connected successfully")
	return nil
}
