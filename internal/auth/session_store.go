package auth

import (
	"context"
	"errors"
)

type SessionInfo struct {
	UserID  string
	Version int
}

var ErrNotFound = errors.New("not found")

type SessionStore interface {
	Set(
		ctx context.Context,
		sessionID string,
		value *SessionInfo,
	) (string, error)
	Get(ctx context.Context, sessionID string) (*SessionInfo, error)
	Delete(ctx context.Context, sessionID string) error
}
