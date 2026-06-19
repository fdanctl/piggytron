package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, u *user.User) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, name, password_hash, created_at, updated_at)
		 VALUES($1,$2,$3,$4,$5)`,
		u.ID(),
		u.Name(),
		u.PasswordHash(),
		u.CreatedAt(),
		u.UpdatedAt(),
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return user.ErrDuplicate
			}
		}

		return err
	}

	return nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, user.ErrNotFound
		}
		return nil, err
	}

	u := user.Rehydrate(id, name, passwordHash, createdAt, updatedAt)
	return u, err
}
