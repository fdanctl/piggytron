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

type SessionVersionStore struct {
	client *redis.Client
}

func NewSessionVersionStore(client *redis.Client) *SessionVersionStore {
	return &SessionVersionStore{
		client: client,
	}
}

func (s *SessionVersionStore) Set(ctx context.Context, userID string, version int) error {
	key := fmt.Sprint(sessionVersionPrefix, userID)

	return s.client.Set(ctx, key, version, 0).Err()
}

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

func (s *SessionVersionStore) Increment(ctx context.Context, userID string) error {
	_, err := s.client.Incr(ctx, fmt.Sprint(sessionVersionPrefix, userID)).Result()
	return err
}

func (s *SessionVersionStore) Delete(ctx context.Context, userID string) error {
	_, err := s.client.Del(ctx, fmt.Sprint(sessionVersionPrefix, userID)).Result()
	return err
}
