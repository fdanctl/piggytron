package middleware

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fdanctl/piggytron/internal/auth"
)

// ErrInvalidSession is returned by SessionInfoFromCtx when the request
// context carries no authenticated session.
var ErrInvalidSession = errors.New("invalid session")

// SessionInfoFromCtx extracts the authenticated session from the request
// context; it errors with ErrInvalidSession when absent.
func SessionInfoFromCtx(ctx context.Context) (*auth.SessionInfo, error) {
	if v, ok := ctx.Value(UserKey).(*auth.SessionInfo); ok {
		return v, nil
	}
	return nil, ErrInvalidSession
}

// LoggerFromContext returns the request-scoped logger, or slog.Default()
// when the context carries none.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

// RequestIDFromContext returns the request id stored by RequestIDMiddleware,
// or "" when absent.
func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
