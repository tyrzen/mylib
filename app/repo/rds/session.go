package rds

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
	"github.com/delveper/mylib/app/exc"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

type Session struct{ *redis.Client }

func NewSession(client *redis.Client) Session {
	return Session{client}
}

func (t Session) Add(ctx context.Context, token ent.Token) error {
	if err := t.Set(ctx, token.ID, token.UID, token.Expiration).Err(); err != nil {
		return fmt.Errorf("recording session: %w", err)
	}

	return nil
}

func (t Session) GetByID(ctx context.Context, id string) (ent.Token, error) {
	uid, err := t.Get(ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return ent.Token{}, fmt.Errorf("nil record: %w", exc.ErrTokenNotFound)
	}

	if err != nil {
		return ent.Token{}, fmt.Errorf("error fetching record: %w", exc.ErrUnexpected)
	}

	token := ent.NewToken(id, uid, 0)

	return token, nil
}
