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
	UserId         string `redis:"user_id"`
	SessionVersion uint   `redis:"session_version"`
}

func (ss *SessionStore) Set(ctx context.Context, value *SessionInfo) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionId := hex.EncodeToString(b)

	key := fmt.Sprint(sessionPrefix, sessionId)
	err = ss.client.HSet(ctx, key, *value).Err()
	ss.client.Expire(ctx, key, time.Hour*24)

	return sessionId, err
}

func (ss *SessionStore) Get(ctx context.Context, sessionId string) *SessionInfo {
	cmd := ss.client.HGetAll(ctx, fmt.Sprint(sessionPrefix, sessionId))
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

func (ss *SessionStore) Remove(ctx context.Context, sessionId string) error {
	_, err := ss.client.Del(ctx, sessionId).Result()
	return err
}
