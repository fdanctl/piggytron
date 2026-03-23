package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{
		db: db,
	}
}

type TransactionDto struct {
	id     transaction.ID
	userId transaction.ID

	ttype transaction.Ttype

	fromAccountId *transaction.ID
	toAccountId   *transaction.ID

	incomeCategoryId  *transaction.ID
	expenseCategoryId *transaction.ID

	amount      uint
	description string
	date        time.Time
	createdAt   time.Time
}

// save

func (r *TransactionRepository) FindById(
	ctx context.Context,
	id transaction.ID,
) (*transaction.Transaction, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE id = $1`,
		id,
	)

	var dto TransactionDto
	err := row.Scan(
		&dto.id,
		&dto.userId,
		&dto.ttype,
		&dto.fromAccountId,
		&dto.toAccountId,
		&dto.incomeCategoryId,
		&dto.expenseCategoryId,
		&dto.amount,
		&dto.description,
		&dto.date,
		&dto.createdAt,
	)
	if err != nil {
		return nil, err
	}

	transaction := transaction.Rehydrate(
		dto.id,
		dto.userId,
		dto.ttype,
		dto.fromAccountId,
		dto.toAccountId,
		dto.incomeCategoryId,
		dto.expenseCategoryId,
		dto.amount,
		dto.description,
		dto.date,
		dto.createdAt,
	)
	return transaction, nil
}

func (r *TransactionRepository) FindAllByUser(
	ctx context.Context,
	uid transaction.ID,
) ([]*transaction.Transaction, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE user_id = $1
		 ORDER BY date DESC`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	for rows.Next() {
		var dto TransactionDto
		err := rows.Scan(
			&dto.id,
			&dto.userId,
			&dto.ttype,
			&dto.fromAccountId,
			&dto.toAccountId,
			&dto.incomeCategoryId,
			&dto.expenseCategoryId,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := transaction.Rehydrate(
			dto.id,
			dto.userId,
			dto.ttype,
			dto.fromAccountId,
			dto.toAccountId,
			dto.incomeCategoryId,
			dto.expenseCategoryId,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) FindAllByCategory(
	ctx context.Context,
	cid transaction.ID,
) ([]*transaction.Transaction, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE income_category_id = $1 OR expense_category_id = $1
		 ORDER BY date DESC`,
		cid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	for rows.Next() {
		var dto TransactionDto
		err := rows.Scan(
			&dto.id,
			&dto.userId,
			&dto.ttype,
			&dto.fromAccountId,
			&dto.toAccountId,
			&dto.incomeCategoryId,
			&dto.expenseCategoryId,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := transaction.Rehydrate(
			dto.id,
			dto.userId,
			dto.ttype,
			dto.fromAccountId,
			dto.toAccountId,
			dto.incomeCategoryId,
			dto.expenseCategoryId,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) FindWithFilters(
	ctx context.Context,
	uid transaction.ID,
	filters *transaction.Filters,
	limit, offset uint,
) ([]*transaction.Transaction, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at
		 FROM transactions
		 WHERE user_id = $1
	     AND ($2::TEXT[] IS NULL OR type = ANY($2))
	     AND ($3::UUID[] IS NULL OR from_account_id = ANY($3) OR to_account_id = ANY($3))
	     AND ($4::UUID[] IS NULL OR income_category_id = ANY($4) OR expense_category_id = ANY($4))
	     AND ($5::BIGINT IS NULL OR amount >= $5)
	     AND ($6::BIGINT IS NULL OR amount <= $6)
		 ORDER BY date DESC
		 LIMIT $7
		 OFFSET $8`,
		uid,
		filters.Ttypes,
		filters.AccountIds,
		filters.CategoryIds,
		filters.MinAmount,
		filters.MaxAmount,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*transaction.Transaction

	for rows.Next() {
		var dto TransactionDto
		err := rows.Scan(
			&dto.id,
			&dto.userId,
			&dto.ttype,
			&dto.fromAccountId,
			&dto.toAccountId,
			&dto.incomeCategoryId,
			&dto.expenseCategoryId,
			&dto.amount,
			&dto.description,
			&dto.date,
			&dto.createdAt,
		)
		if err != nil {
			return nil, err
		}
		transaction := transaction.Rehydrate(
			dto.id,
			dto.userId,
			dto.ttype,
			dto.fromAccountId,
			dto.toAccountId,
			dto.incomeCategoryId,
			dto.expenseCategoryId,
			dto.amount,
			dto.description,
			dto.date,
			dto.createdAt,
		)
		transactions = append(transactions, transaction)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionRepository) CountFilteredResults(
	ctx context.Context,
	uid transaction.ID,
	filters *transaction.Filters,
) (uint, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT COUNT(*)
		 FROM transactions
		 WHERE user_id = $1
	     AND ($2::TEXT[] IS NULL OR type = ANY($2))
	     AND ($3::UUID[] IS NULL OR from_account_id = ANY($3) OR to_account_id = ANY($3))
	     AND ($4::UUID[] IS NULL OR income_category_id = ANY($4) OR expense_category_id = ANY($4))
	     AND ($5::BIGINT IS NULL OR amount >= $5)
	     AND ($6::BIGINT IS NULL OR amount <= $6)`,
		uid,
		filters.Ttypes,
		filters.AccountIds,
		filters.CategoryIds,
		filters.MinAmount,
		filters.MaxAmount,
	)
	var count uint
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}
	fmt.Println(count)
	return count, nil
}
