package rds

import (
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func Connect() *redis.Client {
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

	return client
}
