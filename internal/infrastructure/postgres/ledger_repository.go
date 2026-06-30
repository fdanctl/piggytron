package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/ledger"
)

type LedgerRepository struct {
	db DBTX
}

func NewLedgerRepository(db DBTX) *LedgerRepository {
	return &LedgerRepository{
		db: db,
	}
}

type LedgerEntryDto struct {
	id     ledger.ID
	userID ledger.ID

	ttype ledger.Type

	fromAccountID *ledger.ID
	toAccountID   *ledger.ID

	incomeCategoryID  *ledger.ID
	expenseCategoryID *ledger.ID

	amount      int
	description string
	date        time.Time
	createdAt   time.Time
}

func (r *LedgerRepository) Create(ctx context.Context, t *ledger.Entry) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO ledger (
		    id,
		    user_id,
		    type,
		    from_account_id,
			to_account_id,
		    income_category_id,
		    expense_category_id,
		    amount,
		    description,
		    date,
		    created_at
		 )
	 	 VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		t.ID(),
		t.UserID(),
		t.Type(),
		t.FromAccountID(),
		t.ToAccountID(),
		t.IncomeCategoryID(),
		t.ExpenseCategoryID(),
		t.Amount(),
		t.Description(),
		t.Date(),
		t.CreatedAt(),
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *LedgerRepository) UpdateMany(
	ctx context.Context,
	tt []*ledger.Entry,
) error {
	if len(tt) == 0 {
		return nil
	}
	var (
		args         []any
		valueStrings []string
	)

	const vn = 10 // number of field in VALUES
	for i, t := range tt {
		n := i*vn + 1

		valueStrings = append(
			valueStrings,
			fmt.Sprintf(
				`(
				$%d::uuid,
				$%d,
				$%d::uuid,
				$%d::uuid,
				$%d::uuid,
				$%d::uuid,
				$%d::bigint,
				$%d,
				$%d::timestamp,
				$%d::timestamp
				)`,
				n,
				n+1,
				n+2,
				n+3,
				n+4,
				n+5,
				n+6,
				n+7,
				n+8,
				n+9,
			),
		)

		args = append(
			args,
			t.ID(),
			t.Type(),
			t.FromAccountID(),
			t.ToAccountID(),
			t.IncomeCategoryID(),
			t.ExpenseCategoryID(),
			t.Amount(),
			t.Description(),
			t.Date(),
			t.CreatedAt(),
		)
	}

	query := fmt.Sprintf(`
		UPDATE ledger AS t
		SET
			type = v.type,
			from_account_id = v.from_account_id,
			to_account_id = v.to_account_id,
			income_category_id = v.income_category_id,
			expense_category_id = v.expense_category_id,
			amount = v.amount,
			description = v.description,
			date = v.date,
			created_at = v.created_at
		FROM (
			VALUES %s
		) AS v(
			id,
			type,
		    from_account_id,
		    to_account_id,
		    income_category_id,
		    expense_category_id,
		    amount,
		    description,
		    date,
		    created_at
		)
		WHERE t.id = v.id
	`, strings.Join(valueStrings, ","))

	_, err := r.db.ExecContext(ctx, query, args...)

	return err
}

func (r *LedgerRepository) Delete(ctx context.Context, id ledger.ID) error {
	_, err := r.db.ExecContext(
		ctx,
		`DELETE FROM ledger WHERE id = $1`,
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *LedgerRepository) FindByID(
	ctx context.Context,
	id ledger.ID,
) (*ledger.Entry, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM ledger
		 WHERE id = $1`,
		id,
	)

	var dto LedgerEntryDto
	err := row.Scan(
		&dto.id,
		&dto.userID,
		&dto.ttype,
		&dto.fromAccountID,
		&dto.toAccountID,
		&dto.incomeCategoryID,
		&dto.expenseCategoryID,
		&dto.amount,
		&dto.description,
		&dto.date,
		&dto.createdAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ledger.ErrNotFound
		}
		return nil, err
	}

	transaction := ledger.Rehydrate(
		dto.id,
		dto.userID,
		dto.ttype,
		dto.fromAccountID,
		dto.toAccountID,
		dto.incomeCategoryID,
		dto.expenseCategoryID,
		dto.amount,
		dto.description,
		dto.date,
		dto.createdAt,
	)
	return transaction, nil
}

func (r *LedgerRepository) FindAllByUser(
	ctx context.Context,
	uid ledger.ID,
) ([]*ledger.Entry, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM ledger
		 WHERE user_id = $1
		 ORDER BY date DESC, created_at DESC`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*ledger.Entry

	for rows.Next() {
		var dto LedgerEntryDto
		err := rows.Scan(
			&dto.id,
			&dto.userID,
			&dto.ttype,
			&dto.fromAccountID,
			&dto.toAccountID,
			&dto.incomeCategoryID,
			&dto.expenseCategoryID,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := ledger.Rehydrate(
			dto.id,
			dto.userID,
			dto.ttype,
			dto.fromAccountID,
			dto.toAccountID,
			dto.incomeCategoryID,
			dto.expenseCategoryID,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		entries = append(entries, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *LedgerRepository) FindAllByAccount(
	ctx context.Context,
	aid ledger.ID,
) ([]*ledger.Entry, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM ledger
		 WHERE from_account_id = $1 OR to_account_id = $1
		 ORDER BY date DESC, created_at DESC`,
		aid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*ledger.Entry

	for rows.Next() {
		var dto LedgerEntryDto
		err := rows.Scan(
			&dto.id,
			&dto.userID,
			&dto.ttype,
			&dto.fromAccountID,
			&dto.toAccountID,
			&dto.incomeCategoryID,
			&dto.expenseCategoryID,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := ledger.Rehydrate(
			dto.id,
			dto.userID,
			dto.ttype,
			dto.fromAccountID,
			dto.toAccountID,
			dto.incomeCategoryID,
			dto.expenseCategoryID,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		entries = append(entries, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (r *LedgerRepository) FindAllByCategory(
	ctx context.Context,
	cid ledger.ID,
) ([]*ledger.Entry, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM ledger
		 WHERE income_category_id = $1 OR expense_category_id = $1
		 ORDER BY date DESC, created_at DESC`,
		cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []*ledger.Entry

	for rows.Next() {
		var dto LedgerEntryDto
		err := rows.Scan(
			&dto.id,
			&dto.userID,
			&dto.ttype,
			&dto.fromAccountID,
			&dto.toAccountID,
			&dto.incomeCategoryID,
			&dto.expenseCategoryID,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := ledger.Rehydrate(
			dto.id,
			dto.userID,
			dto.ttype,
			dto.fromAccountID,
			dto.toAccountID,
			dto.incomeCategoryID,
			dto.expenseCategoryID,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		entries = append(entries, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}
