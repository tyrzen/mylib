package rds

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Token struct {
	client *redis.Client
}

func NewToken(client *redis.Client) Token {
	return Token{client}
}

func (s Token) Create(ctx context.Context, token ent.Token) error {
	if err := s.client.Set(ctx, token.ID, token.UID, token.Expiry).Err(); err != nil {
		return fmt.Errorf("recording session: %w", err)
	}

	return nil
}

func (s Token) Find(ctx context.Context, token ent.Token) (ent.Token, error) {
	uid, err := s.client.Get(ctx, token.ID).Result()
	if errors.Is(err, redis.Nil) {
		return ent.Token{}, fmt.Errorf("nil record: %w", exc.ErrTokenNotFound)
	}

	if err != nil {
		return ent.Token{}, fmt.Errorf("error fetching record: %w", exc.ErrUnexpected)
	}

	token.UID = uid

	return token, nil
}

func (s Token) Destroy(ctx context.Context, token ent.Token) error {
	_, err := s.client.Del(ctx, token.ID).Result()

	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("nil record: %w", exc.ErrTokenNotFound)
	}

	if err != nil {
		return fmt.Errorf("error fetching record: %w", exc.ErrUnexpected)
	}

	return nil
}
