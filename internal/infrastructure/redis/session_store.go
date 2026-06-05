package redis

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

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

type SessionInfo struct {
	UserID         string `redis:"user_id"`
	SessionVersion uint   `redis:"session_version"`
}

func (s *SessionStore) Set(ctx context.Context, value *SessionInfo) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(b)

	key := fmt.Sprint(sessionPrefix, sessionID)
	err = s.client.HSet(ctx, key, *value).Err()
	s.client.Expire(ctx, key, time.Hour*24)

	return sessionID, err
}

func (s *SessionStore) Get(ctx context.Context, sessionID string) *SessionInfo {
	cmd := s.client.HGetAll(ctx, fmt.Sprint(sessionPrefix, sessionID))
	m, err := cmd.Result()
	if err != nil {
		return nil
	}

	if len(m) == 0 {
		return nil
	}
	var res SessionInfo
	if err := cmd.Scan(&res); err != nil {
		return nil
	}
	return &res
}

func (s *SessionStore) Remove(ctx context.Context, sessionID string) error {
	_, err := s.client.Del(ctx, fmt.Sprint(sessionPrefix, sessionID)).Result()
	return err
}
