package auth

import (
	"context"
	"errors"
)

// SessionInfo is the payload stored with each session id.
type SessionInfo struct {
	UserID  string
	Version int // per-user session version at issue time
}

// ErrNotFound is returned by stores when a key has no value.
var ErrNotFound = errors.New("not found")

// SessionStore persists session ids to their SessionInfo payloads;
// implemented by the redis package.
type SessionStore interface {
	Set(
		ctx context.Context,
		sessionID string,
		value *SessionInfo,
	) (string, error)
	Get(ctx context.Context, sessionID string) (*SessionInfo, error)
	Delete(ctx context.Context, sessionID string) error
}
