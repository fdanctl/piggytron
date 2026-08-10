package auth

import "context"

// SessionVersionStore is the contract for a per-user session version store.
// Bumping the stored version invalidates previously issued sessions.
type SessionVersionStore interface {
	Set(ctx context.Context, userID string, version int) error
	Get(ctx context.Context, userID string) (int, error)
	Increment(ctx context.Context, userID string) error
	Delete(ctx context.Context, userID string) error
}
