package auth

import "context"

type SessionVersionStore interface {
	Set(ctx context.Context, userID string, version int) error
	Get(ctx context.Context, userID string) (int, error)
	Increment(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
}
