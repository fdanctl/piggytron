package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

type SessionManager struct {
	sessionStore        SessionStore
	sessionVersionStore SessionVersionStore
}

func NewSessionManager(ss SessionStore, svs SessionVersionStore) *SessionManager {
	return &SessionManager{
		sessionStore:        ss,
		sessionVersionStore: svs,
	}
}

func (m *SessionManager) CreateSession(
	ctx context.Context,
	userID string,
) (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	sessionID := hex.EncodeToString(b)

	version, err := m.sessionVersionStore.Get(ctx, userID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			version = 1
			err = m.sessionVersionStore.Set(ctx, userID, 1)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}

	value := SessionInfo{
		UserID:  userID,
		Version: version,
	}

	return m.sessionStore.Set(ctx, sessionID, &value)
}

func (m *SessionManager) GetSession(
	ctx context.Context,
	sessionID string,
) (*SessionInfo, error) {
	return m.sessionStore.Get(ctx, sessionID)
}

func (m *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return m.sessionStore.Delete(ctx, sessionID)
}

func (m *SessionManager) RevokeAllSessions(ctx context.Context, userID string) error {
	return m.sessionVersionStore.Increment(ctx, userID)
}

func (m *SessionManager) ValidateSession(
	ctx context.Context,
	session *SessionInfo,
) (bool, error) {
	version, err := m.sessionVersionStore.Get(ctx, session.UserID)
	if err != nil {
		return false, err
	}
	if version != session.Version {
		return false, nil
	}
	return true, nil
}
