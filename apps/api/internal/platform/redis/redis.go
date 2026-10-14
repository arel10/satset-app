package redis

import (
	"context"
	"fmt"

	goredis "github.com/redis/go-redis/v9"
)

type Client struct {
	*goredis.Client
}

func Connect(addr, password string, db int) (*Client, error) {
	client := goredis.NewClient(&goredis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	return &Client{Client: client}, nil
}

func (c *Client) Ping(ctx context.Context) error {
	if c == nil || c.Client == nil {
		return fmt.Errorf("redis client not initialized")
	}
	return c.Client.Ping(ctx).Err()
}
