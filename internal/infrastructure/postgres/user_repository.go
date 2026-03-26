package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Save(ctx context.Context, user *user.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, name, password_hash, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5)`,
		user.ID(),
		user.Name(),
		user.PasswordHash(),
		user.CreatedAt(),
		user.UpdatedAt(),
	)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id user.ID) (*user.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, password_hash, created_at, updated_at
		 FROM users
		 WHERE id = $1`,
		id,
	)

	var (
		name         string
		passwordHash string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(
		&id,
		&name,
		&passwordHash,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u := user.Rehydrate(id, name, passwordHash, createdAt, updatedAt)
	return u, err
}

func (r *UserRepository) FindByName(ctx context.Context, name string) (*user.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, password_hash, created_at, updated_at
		 FROM users
		 WHERE name = $1`,
		name,
	)

	var (
		id           user.ID
		passwordHash string
		createdAt    time.Time
		updatedAt    time.Time
	)

	err := row.Scan(
		&id,
		&name,
		&passwordHash,
		&createdAt,
		&updatedAt,
	)
	if err != nil {
		return nil, err
	}

	u := user.Rehydrate(id, name, passwordHash, createdAt, updatedAt)
	return u, err
}
