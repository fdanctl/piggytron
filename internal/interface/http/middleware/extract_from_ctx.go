package middleware

import (
	"context"
	"errors"
	"log/slog"

	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

var ErrInvalidSession = errors.New("invalid session")

func SessionInfoFromCtx(ctx context.Context) (*rdb.SessionInfo, error) {
	if v, ok := ctx.Value(UserKey).(*rdb.SessionInfo); ok {
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
