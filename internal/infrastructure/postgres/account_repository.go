package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/lib/pq"
)

// AccountRepository persists account aggregates via raw SQL.
type AccountRepository struct {
	db DBTX
}

// NewAccountRepository builds an account repository over a DBTX (sql.DB or
// sql.Tx).
func NewAccountRepository(db DBTX) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// accountDto is the database row shape of the accounts table.
type accountDto struct {
	ID     account.ID
	UserID account.ID
	Type   account.AccountType
	Name   string

	Status       account.AccountStatus
	TargetAmount *int
	StartDate    *time.Time
	TargetDate   *time.Time
	CategoryID   *account.ID

	Currency  string
	ClosedAt  *time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Create inserts the account, mapping unique-name violations to
// account.ErrDuplicate.
func (r *AccountRepository) Create(ctx context.Context, a *account.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, user_id, type, name, status, currency, target_amount, start_date, target_date, category_id, closed_at, created_at, updated_at)
	 	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`,
		a.ID(),
		a.UserID(),
		a.Type(),
		a.Name(),
		a.Status(),
		a.Currency(),
		a.TargetAmount(),
		a.StartDate(),
		a.TargetDate(),
		a.CategoryID(),
		a.ClosedAt(),
		a.CreatedAt(),
		a.UpdatedAt(),
	)
	if err != nil {
		var pqErr *pq.Error

		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case "23505":
				return account.ErrDuplicate
			}
		}

		return err
	}

	return nil
}

// Update persists the mutable fields of an account.
func (r *AccountRepository) Update(ctx context.Context, a *account.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE accounts
		SET
			name = $2,
			status = $3,
			currency = $4,
			target_amount = $5,
			start_date = $6,
			target_date = $7,
			category_id = $8,
			closed_at = $9,
			updated_at = $10
		WHERE id = $1`,
		a.ID(),
		a.Name(),
		a.Status(),
		a.Currency(),
		a.TargetAmount(),
		a.StartDate(),
		a.TargetDate(),
		a.CategoryID(),
		a.ClosedAt(),
		a.UpdatedAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes matching id account
func (r *AccountRepository) Delete(ctx context.Context, id account.ID) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM accounts WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

// FindByID loads one account, mapping missing rows to account.ErrNotFound.
func (r *AccountRepository) FindByID(ctx context.Context, id account.ID) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, name, status, currency, target_amount, start_date, target_date, category_id, closed_at, created_at, updated_at
		 FROM accounts
		 WHERE id = $1`,
		id,
	)

	var b accountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Type,
		&b.Name,
		&b.Status,
		&b.Currency,
		&b.TargetAmount,
		&b.StartDate,
		&b.TargetDate,
		&b.CategoryID,
		&b.ClosedAt,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, account.ErrNotFound
		}
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		b.Type,
		b.Name,
		b.Status,
		b.TargetAmount,
		b.StartDate,
		b.TargetDate,
		b.CategoryID,
		b.Currency,
		b.ClosedAt,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

// FindAllByUser returns every account of the user.
func (r *AccountRepository) FindAllByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, status, currency, target_amount, start_date, target_date, category_id, closed_at, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto accountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.Status,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.StartDate,
			&dto.TargetDate,
			&dto.CategoryID,
			&dto.ClosedAt,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			dto.Type,
			dto.Name,
			dto.Status,
			dto.TargetAmount,
			dto.StartDate,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
			dto.ClosedAt,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindAllBanksByUser returns the bank accounts of the user.
func (r *AccountRepository) FindAllBanksByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, status, currency, target_amount, start_date, target_date, category_id, closed_at, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND type != 'goal'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto accountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.Status,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.StartDate,
			&dto.TargetDate,
			&dto.CategoryID,
			&dto.ClosedAt,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			dto.Type,
			dto.Name,
			dto.Status,
			dto.TargetAmount,
			dto.StartDate,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
			dto.ClosedAt,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}

// FindAllGoalsByUser returns the goal accounts of the user.
func (r *AccountRepository) FindAllGoalsByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, status, currency, target_amount, start_date, target_date, category_id, closed_at, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND type = 'goal'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto accountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.Status,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.StartDate,
			&dto.TargetDate,
			&dto.CategoryID,
			&dto.ClosedAt,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			dto.Type,
			dto.Name,
			dto.Status,
			dto.TargetAmount,
			dto.StartDate,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
			dto.ClosedAt,
			dto.CreatedAt,
			dto.UpdatedAt,
		)
		accounts = append(accounts, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return accounts, nil
}
