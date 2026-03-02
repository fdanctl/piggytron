package user

import "context"

type Repository interface {
	Save(ctx context.Context, user *User) error
	FindById(ctx context.Context, id ID) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
}
