package config

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx   = context.Background()
	Redis *redis.Client
)

func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_ADDR"), // force IPv4
	})

	// test connection
	_, err := Redis.Ping(Ctx).Result()
	return err
}
