package config

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var (
	Ctx   = context.Background()
	Redis *redis.Client
)

func InitRedis() error {
	Redis = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379", // force IPv4
	})

	// test connection
	_, err := Redis.Ping(Ctx).Result()
	return err
}
