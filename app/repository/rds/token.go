package rds

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Token struct {
	client *redis.Client
}

func NewToken(client *redis.Client) *Token {
	return &Token{client}
}

func (t Token) Create(ctx context.Context, token models.Token) error {
	if err := t.client.Set(ctx, token.ID, token.UID, token.Expiry).Err(); err != nil {
		return fmt.Errorf("recording session: %w", err)
	}

	return nil
}

func (t Token) Find(ctx context.Context, token models.Token) (models.Token, error) {
	uid, err := t.client.Get(ctx, token.ID).Result()
	if errors.Is(err, redis.Nil) {
		return models.Token{}, fmt.Errorf("nil record: %w", exceptions.ErrTokenNotFound)
	}

	if err != nil {
		return models.Token{}, fmt.Errorf("error fetching record: %w", exceptions.ErrUnexpected)
	}

	token.UID = uid

	return token, nil
}

func (t Token) Destroy(ctx context.Context, token models.Token) error {
	_, err := t.client.Del(ctx, token.ID).Result()

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("nil record: %w", exceptions.ErrTokenNotFound)
	}

	if err != nil {
		return fmt.Errorf("error fetching record: %w", exceptions.ErrUnexpected)
	}

	return nil
}
