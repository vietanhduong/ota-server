package redis

import (
	"github.com/go-redis/redis"
)

type Client struct {
	*redis.Client
}

func InitializeConnection(cfg Config) (*Client, error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}

	return &Client{redisClient}, nil
}
