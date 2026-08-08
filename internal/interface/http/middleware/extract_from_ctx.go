package middleware

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fdanctl/piggytron/internal/auth"
)

var ErrInvalidSession = errors.New("invalid session")

func SessionInfoFromCtx(ctx context.Context) (*auth.SessionInfo, error) {
	if v, ok := ctx.Value(UserKey).(*auth.SessionInfo); ok {
		return v, nil
	}
	return nil, ErrInvalidSession
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(LoggerKey).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}

func RequestIDFromContext(ctx context.Context) string {
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
