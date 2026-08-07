package user

import "context"

type Repository interface {
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	UpdateName(ctx context.Context, id ID, name string) error
	FindByID(ctx context.Context, id ID) (*User, error)
	FindByName(ctx context.Context, name string) (*User, error)
}
