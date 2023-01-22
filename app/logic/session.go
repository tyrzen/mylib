package logic

import (
	"context"
	"fmt"

	"github.com/delveper/mylib/app/ent"
)

type Session struct {
	SessionRepository
}

func (s Session) Create(ctx context.Context, token ent.Token) error {
	if err := s.Add(ctx, token); err != nil {
		return fmt.Errorf("session error: %w", err)
	}

	return nil
}

func (s Session) Find(ctx context.Context, id string) (ent.Token, error) {
	token, err := s.GetByID(ctx, id)
	if err != nil {
		return ent.Token{}, fmt.Errorf("session error: %w", err)
	}

	return token, nil
}
