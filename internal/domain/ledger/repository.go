package ledger

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, entry *Entry) error
	UpdateMany(ctx context.Context, entries []*Entry) error
	Delete(ctx context.Context, id ID) error
	FindByID(ctx context.Context, id ID) (*Entry, error)
	FindAllByUser(ctx context.Context, uid ID) ([]*Entry, error)
	FindAllByAccount(ctx context.Context, aid ID) ([]*Entry, error)
	FindAllByCategory(ctx context.Context, cid ID) ([]*Entry, error)
}
