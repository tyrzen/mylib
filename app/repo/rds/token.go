package rds

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Token struct{ *redis.Client }

func NewToken(client *redis.Client) Token {
	return Token{client}
}

func (t Token) Add(ctx context.Context, token ent.Token) error {
	if err := t.Set(ctx, token.ID, token.UID, token.Expiration).Err(); err != nil {
		return fmt.Errorf("error creating session record: %w", err)
	}

	return nil
}

func (t Token) GetByID(ctx context.Context, id string) (ent.Token, error) {
	uid, err := t.Get(ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return ent.Token{}, fmt.Errorf("nil record: %w", exc.ErrTokenNotFound)
	}

	if err != nil {
		return ent.Token{}, fmt.Errorf("failed fetching record: %w", exc.ErrUnexpected)
	}

	token := ent.Token{
		ID:  id,
		UID: uid,
	}

	return token, nil
}
