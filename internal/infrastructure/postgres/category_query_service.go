package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
)

type CategoryQueryService struct {
	db *sql.DB
}

func NewCategoryQueryService(db *sql.DB) *CategoryQueryService {
	return &CategoryQueryService{
		db: db,
	}
}

func (s *CategoryQueryService) FindAllCategories(
	ctx context.Context,
	uid string,
) ([]query.CategoryNameDTO, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name
		 FROM income_categories
		 WHERE user_id = $1
		 UNION
		 SELECT id, name
		 FROM expense_categories
		 WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryNameDTO

	for rows.Next() {
		var c query.CategoryNameDTO
		if err := rows.Scan(
			&c.ID,
			&c.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *CategoryQueryService) FindCategoriesIDIncludes(
	ctx context.Context,
	ids []string,
) ([]query.CategoryNameDTO, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	queryStr := fmt.Sprintf(
		`SELECT id, name
		 FROM income_categories
		 WHERE id IN (%s)
		 UNION
		 SELECT id, name
		 FROM expense_categories
		 WHERE id IN (%s)`,
		strings.Join(placeholders, ","),
		strings.Join(placeholders, ","),
	)

	rows, err := s.db.QueryContext(
		ctx,
		queryStr,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []query.CategoryNameDTO

	for rows.Next() {
		var c query.CategoryNameDTO
		if err := rows.Scan(
			&c.ID,
			&c.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (s *CategoryQueryService) GetExpenseCategoriesBudgetSpent(
	ctx context.Context,
	uid string,
	minDate, maxDate time.Time,
) ([]query.ExpenseCategoryBudgetSpent, error) {
	rows, err := s.db.QueryContext(
		ctx,
		`
        SELECT
          c.id as cid,
          COALESCE(b.id, '00000000-0000-0000-0000-000000000000') as bid,
          c.type,
          c.name,
          COALESCE(b.amount, 0),
          COALESCE(SUM(t.amount), 0)
        FROM
          expense_categories c
          LEFT JOIN transactions t ON c.id = t.expense_category_id
          AND t.date >= $1
          AND t.date < $2
          LEFT JOIN monthly_budgets b ON c.id = b.category_id
          AND b.month >= $1
          AND b.month < $2
        WHERE
          c.user_id = $3
        GROUP BY
          c.id,
          b.id`,
		minDate,
		maxDate,
		uid,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var results []query.ExpenseCategoryBudgetSpent

	for rows.Next() {
		var dto query.ExpenseCategoryBudgetSpent
		if err := rows.Scan(
			&dto.CID,
			&dto.BID,
			&dto.Type,
			&dto.Name,
			&dto.Budgeted,
			&dto.Spent,
		); err != nil {
			return nil, err
		}
		results = append(results, dto)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
