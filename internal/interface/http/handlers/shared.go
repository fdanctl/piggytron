package handlers

import (
	"context"
	"errors"

	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
)

const LIMIT = 30

var ErrInvalidSession = errors.New("invalid session")

func sessionInfoFormCtx(ctx context.Context) (*rdb.SessionInfo, error) {
	v := ctx.Value(middleware.UserKey)
	if v == nil {
		return nil, ErrInvalidSession
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, ErrInvalidSession
	}
	return sessionInfo, nil
}
