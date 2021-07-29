package redis

import (
	"fmt"
	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func InitializeConnection(cfg Config) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprint(cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}

	return &Client{redisClient}, nil
}
