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

type SessionStore struct {
	client *redis.Client
}

func NewSessionStore(client *redis.Client) *SessionStore {
	return &SessionStore{
		client: client,
	}
}

type SessionInfoDTO struct {
	UserID  string `redis:"user_id"`
	Version int    `redis:"session_version"`
}

func (s *SessionStore) Set(
	ctx context.Context,
	sessionID string,
	value *auth.SessionInfo,
) (string, error) {
	dto := SessionInfoDTO{
		UserID:  value.UserID,
		Version: value.Version,
	}
	key := fmt.Sprint(sessionPrefix, sessionID)
	err := s.client.HSet(ctx, key, dto).Err()
	s.client.Expire(ctx, key, time.Hour*24)

	return sessionID, err
}

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
	var res SessionInfoDTO
	if err := cmd.Scan(&res); err != nil {
		return nil, err
	}
	return &auth.SessionInfo{
		UserID:  res.UserID,
		Version: res.Version,
	}, nil
}

func (s *SessionStore) Delete(ctx context.Context, sessionID string) error {
	_, err := s.client.Del(ctx, fmt.Sprint(sessionPrefix, sessionID)).Result()
	return err
}
