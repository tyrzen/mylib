package rds

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func Connect() (*redis.Client, error) {
	addr := fmt.Sprintf("%s:%s",
		os.Getenv("SESSION_HOST"),
		os.Getenv("SESSION_PORT"),
	)

	pwd := os.Getenv("SESSION_PASSWORD")

	cfg := &redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0,
	}

	client := redis.NewClient(cfg)

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("error pinging repo: %w", err)
	}

	return client, nil
}
