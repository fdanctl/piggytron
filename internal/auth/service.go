// Package auth manages user sessions. A session is a random hex id stored in
// Redis together with the user id and a per-user version counter; bumping the
// counter (RevokeAllSessions) invalidates every previously issued session.
package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

// SessionManager creates and validates sessions against the per-user version
// counter. It composes a SessionStore and a SessionVersionStore.
type SessionManager struct {
	sessionStore        SessionStore
	sessionVersionStore SessionVersionStore
}

// NewSessionManager wires a session manager to its stores.
func NewSessionManager(ss SessionStore, svs SessionVersionStore) *SessionManager {
	return &SessionManager{
		sessionStore:        ss,
		sessionVersionStore: svs,
	}
}

// CreateSession issues a new session id for the user, stamping it with the
// user's current session version.
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

// GetSession looks up a session by its id.
func (m *SessionManager) GetSession(
	ctx context.Context,
	sessionID string,
) (*SessionInfo, error) {
	return m.sessionStore.Get(ctx, sessionID)
}

// DeleteSession removes a single session (logout).
func (m *SessionManager) DeleteSession(ctx context.Context, sessionID string) error {
	return m.sessionStore.Delete(ctx, sessionID)
}

// RevokeAllSessions invalidates every session of the user by bumping the
// per-user version counter.
func (m *SessionManager) RevokeAllSessions(ctx context.Context, userID string) error {
	return m.sessionVersionStore.Increment(ctx, userID)
}

// ValidateSession reports whether the session's recorded version still
// matches the user's current version (false after RevokeAllSessions).
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
