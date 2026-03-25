package postgres

import (
	"context"
	"database/sql"

	"github.com/fdanctl/piggytron/internal/query"
)

type TransactionQueryService struct {
	db *sql.DB
}

func NewTransactionQueryService(db *sql.DB) *TransactionQueryService {
	return &TransactionQueryService{
		db: db,
	}
}

func (r *TransactionQueryService) FindFiltered(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
	limit, offset uint,
) ([]query.TransactionDTO, error) {
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
		filters.Types,
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

	var transactions []query.TransactionDTO

	for rows.Next() {
		var dto query.TransactionDTO
		err := rows.Scan(
			&dto.Id,
			&dto.UserId,
			&dto.Type,
			&dto.FromAccountId,
			&dto.ToAccountId,
			&dto.IncomeCategoryId,
			&dto.ExpenseCategoryId,
			&dto.Amount,
			&dto.Description,
			&dto.Date,
			&dto.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func (r *TransactionQueryService) FindFilteredWithCount(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
	limit, offset uint,
) (*query.TransactionsWithTotalCount, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, type, from_account_id, to_account_id, income_category_id, expense_category_id, amount, description, date, created_at, COUNT(*) OVER() AS total_count
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
		filters.Types,
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

	var transactions []query.TransactionDTO
	var totalCount int

	for rows.Next() {
		var dto query.TransactionDTO
		var count int
		err := rows.Scan(
			&dto.Id,
			&dto.UserId,
			&dto.Type,
			&dto.FromAccountId,
			&dto.ToAccountId,
			&dto.IncomeCategoryId,
			&dto.ExpenseCategoryId,
			&dto.Amount,
			&dto.Description,
			&dto.Date,
			&dto.CreatedAt,
			&count,
		)
		if err != nil {
			return nil, err
		}
		totalCount = count
		transactions = append(transactions, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &query.TransactionsWithTotalCount{
		Data:  transactions,
		Total: totalCount,
	}, nil
}

func (r *TransactionQueryService) CountFilteredResults(
	ctx context.Context,
	uid string,
	filters *query.TransactionFilters,
) (int, error) {
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
		filters.Types,
		filters.AccountIds,
		filters.CategoryIds,
		filters.MinAmount,
		filters.MaxAmount,
	)
	var count int
	err := row.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
