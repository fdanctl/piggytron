// Package redis implements the auth session stores (SessionStore and
// SessionVersionStore) on top of go-redis.
package redis

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/auth"
	"github.com/redis/go-redis/v9"
)

const (
	sessionPrefix = "session:"
)

// SessionStore persists sessions as Redis hashes keyed "session:<id>",
// expiring after 24 hours.
type SessionStore struct {
	client *redis.Client
}

// NewSessionStore builds the store over a Redis client.
func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{
		client: client,
	}
}

// sessionInfoDTO is the redis hash shape of a session.
type sessionInfoDTO struct {
	UserID  string `redis:"user_id"`
	Version int    `redis:"session_version"`
}

// Set stores a session with a 24-hour TTL and returns the session id.
func (s *SessionStore) Set(
	ctx context.Context,
	sessionID string,
	value *auth.SessionInfo,
) (string, error) {
	dto := sessionInfoDTO{
		UserID:  value.UserID,
		Version: value.Version,
	}
	key := fmt.Sprint(sessionPrefix, sessionID)
	err := s.client.HSet(ctx, key, dto).Err()
	s.client.Expire(ctx, key, time.Hour*24)

	return sessionID, err
}

// Get loads a session, mapping missing keys to auth.ErrNotFound.
func (s *SessionStore) Get(ctx context.Context, sessionID string) (*auth.SessionInfo, error) {
	cmd := s.client.HGetAll(ctx, fmt.Sprint(sessionPrefix, sessionID))
	m, err := cmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, auth.ErrNotFound
		}
		return nil, err
	}

	if len(m) == 0 {
		return nil, auth.ErrNotFound
	}
	var res sessionInfoDTO
	if err := cmd.Scan(&res); err != nil {
		return nil, err
	}
	return &auth.SessionInfo{
		UserID:  res.UserID,
		Version: res.Version,
	}, nil
}

// Delete removes a session (logout).
func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	_, err := s.client.Del(ctx, fmt.Sprint(sessionPrefix, sessionID)).Result()
	return err
}
