package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/auth"
	"github.com/redis/go-redis/v9"
)

const (
	sessionVersionPrefix = "session_version:"
)

// SessionVersionStore keeps the per-user session version counter keyed
// "session_version:<userID>"; bumping it invalidates all previously issued
// sessions.
type SessionVersionStore struct {
	client *redis.Client
}

// NewSessionVersionStore builds the store over a Redis client.
func NewSessionVersionStore(client *redis.Client) *SessionVersionStore {
	return &SessionVersionStore{
		client: client,
	}
}

// Set stores the current version for a user.
func (s *SessionVersionStore) Set(ctx context.Context, userID string, version int) error {
	key := fmt.Sprint(sessionVersionPrefix, userID)

	return s.client.Set(ctx, key, version, 0).Err()
}

// Get loads the current version, mapping missing keys to auth.ErrNotFound.
func (s *SessionVersionStore) Get(ctx context.Context, userID string) (int, error) {
	cmd := s.client.Get(ctx, fmt.Sprint(sessionVersionPrefix, userID))
	m, err := cmd.Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return 0, auth.ErrNotFound
		}
		return 0, err
	}
	if len(m) == 0 {
		return 0, auth.ErrNotFound
	}
	var res int
	if err := cmd.Scan(&res); err != nil {
		return 0, err
	}
	return res, nil
}

// Increment bumps the user's version (revoke-all-sessions).
func (s *SessionVersionStore) Increment(ctx context.Context, userID string) error {
	_, err := s.client.Incr(ctx, fmt.Sprint(sessionVersionPrefix, userID)).Result()
	return err
}

// Delete removes the user's version counter.
func (s *SessionVersionStore) Delete(ctx context.Context, userID string) error {
	_, err := s.client.Del(ctx, fmt.Sprint(sessionVersionPrefix, userID)).Result()
	return err
}
