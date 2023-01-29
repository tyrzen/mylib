package rds

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/exceptions"
	"github.com/delveper/mylib/app/models"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Session struct {
	client *redis.Client
}

func NewToken(client *redis.Client) Session {
	return Session{client}
}

func (s Session) Create(ctx context.Context, token models.Token) error {
	if err := s.client.Set(ctx, token.ID, token.UID, token.Expiry).Err(); err != nil {
		return fmt.Errorf("recording session: %w", err)
	}

	return nil
}

func (s Session) Find(ctx context.Context, token models.Token) (models.Token, error) {
	uid, err := s.client.Get(ctx, token.ID).Result()
	if errors.Is(err, redis.Nil) {
		return models.Token{}, fmt.Errorf("nil record: %w", exceptions.ErrTokenNotFound)
	}

	if err != nil {
		return models.Token{}, fmt.Errorf("error fetching record: %w", exceptions.ErrUnexpected)
	}

	token.UID = uid

	return token, nil
}

func (s Session) Destroy(ctx context.Context, token models.Token) error {
	_, err := s.client.Del(ctx, token.ID).Result()

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("nil record: %w", exceptions.ErrTokenNotFound)
	}

	if err != nil {
		return fmt.Errorf("error fetching record: %w", exceptions.ErrUnexpected)
	}

	return nil
}
