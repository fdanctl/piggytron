package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/lib/pq"
)

type AccountRepository struct {
	db DBTX
}

func NewAccountRepository(db DBTX) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

type AccountDto struct {
	ID     account.ID
	UserID account.ID
	Type   account.AccountType
	Name   string

	IsSaving     *bool
	TargetAmount *int
	StartDate    *time.Time
	TargetDate   *time.Time
	CategoryID   *account.ID

	Currency  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (r *AccountRepository) Create(ctx context.Context, a *account.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO accounts (id, user_id, type, name, is_saving, currency, target_amount, start_date, target_date, category_id, created_at, updated_at)
	 	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		a.ID(),
		a.UserID(),
		a.Type(),
		a.Name(),
		a.IsSaving(),
		a.Currency(),
		a.TargetAmount(),
		a.StartDate(),
		a.TargetDate(),
		a.CategoryID(),
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

func (r *AccountRepository) Update(ctx context.Context, a *account.Account) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE accounts
		SET
			name = $2,
			is_saving = $3,
			currency = $4,
			target_amount = $5,
			start_date = $6,
			target_date = $7,
			category_id = $8,
			updated_at = $9
		WHERE id = $1`,
		a.ID(),
		a.Name(),
		a.IsSaving(),
		a.Currency(),
		a.TargetAmount(),
		a.StartDate(),
		a.TargetDate(),
		a.CategoryID(),
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

func (r *AccountRepository) FindByID(ctx context.Context, id account.ID) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, name, is_saving, currency, target_amount, start_date, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE id = $1`,
		id,
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Type,
		&b.Name,
		&b.IsSaving,
		&b.Currency,
		&b.TargetAmount,
		&b.StartDate,
		&b.TargetDate,
		&b.CategoryID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		b.Type,
		b.Name,
		b.IsSaving,
		b.TargetAmount,
		b.StartDate,
		b.TargetDate,
		b.CategoryID,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindBankByNameAndUser(
	ctx context.Context,
	uid account.ID,
	name string,
) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, name, is_saving, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND name = $2 and type = $3`,
		uid,
		name,
		"bank",
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Name,
		&b.IsSaving,
		&b.Currency,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		"bank",
		b.Name,
		b.IsSaving,
		nil,
		nil,
		nil,
		nil,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindGoalByNameAndUser(
	ctx context.Context,
	uid account.ID,
	name string,
) (*account.Account, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, name, is_saving, currency, target_amount, start_date, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND name = $2 and type = $3`,
		uid,
		name,
		"goal",
	)

	var b AccountDto
	err := row.Scan(
		&b.ID,
		&b.UserID,
		&b.Type,
		&b.Name,
		&b.IsSaving,
		&b.Currency,
		&b.TargetAmount,
		&b.StartDate,
		&b.TargetDate,
		&b.CategoryID,
		&b.CreatedAt,
		&b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	account := account.Rehydrate(
		b.ID,
		b.UserID,
		b.Type,
		b.Name,
		b.IsSaving,
		b.TargetAmount,
		b.StartDate,
		b.TargetDate,
		b.CategoryID,
		b.Currency,
		b.CreatedAt,
		b.UpdatedAt,
	)
	return account, err
}

func (r *AccountRepository) FindAllByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, is_saving, currency, target_amount, start_date, target_date, category_id, created_at, updated_at
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
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.IsSaving,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.StartDate,
			&dto.TargetDate,
			&dto.CategoryID,
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
			dto.IsSaving,
			dto.TargetAmount,
			dto.StartDate,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
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

func (r *AccountRepository) FindAllBanksByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, name, is_saving, currency, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 AND type = 'bank'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Name,
			&dto.IsSaving,
			&dto.Currency,
			&dto.CreatedAt,
			&dto.UpdatedAt,
		); err != nil {
			return nil, err
		}
		b := account.Rehydrate(
			dto.ID,
			dto.UserID,
			"bank",
			dto.Name,
			dto.IsSaving,
			nil,
			nil,
			nil,
			nil,
			dto.Currency,
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

func (r *AccountRepository) FindAllGoalsByUser(
	ctx context.Context,
	uid account.ID,
) ([]*account.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, name, is_saving, currency, target_amount, start_date, target_date, category_id, created_at, updated_at
		 FROM accounts
		 WHERE user_id = $1 and type = 'goal'`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*account.Account

	for rows.Next() {
		var dto AccountDto
		if err := rows.Scan(
			&dto.ID,
			&dto.UserID,
			&dto.Type,
			&dto.Name,
			&dto.IsSaving,
			&dto.Currency,
			&dto.TargetAmount,
			&dto.StartDate,
			&dto.TargetDate,
			&dto.CategoryID,
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
			dto.IsSaving,
			dto.TargetAmount,
			dto.StartDate,
			dto.TargetDate,
			dto.CategoryID,
			dto.Currency,
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
