package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type Client struct {
	Redis *redis.Client
}

func NewRedisClient(redisUrl string) (*redis.Client, error) {
	redis := redis.NewClient(&redis.Options{
		Addr:         redisUrl, // "localhost:6379"
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := redis.Ping(ctx).Err(); err != nil {
		panic(err)
	}

	fmt.Println("Redis client initialized")
	return redis, nil
}
