package core_redis

import (
	"context"
	"fmt"
	"net"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	*redis.Client
}

func NewClient(
	ctx context.Context,
	config Config,
) (*Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     net.JoinHostPort(config.Host, config.Port),
		Password: config.Password,
		DB:       config.Database,
	})

	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &Client{client}, nil
}
