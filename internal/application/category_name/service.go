package categoryname

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

type Service struct {
	db *sql.DB
}

type CategoryWithName struct {
	Id   string
	Name string
}

func NewService(db *sql.DB) *Service {
	fmt.Println("db", db)
	return &Service{
		db: db,
	}
}

func (s *Service) GetAllCategories(ctx context.Context) ([]*CategoryWithName, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT id, name
		 FROM income_categories
		 WHERE user_id = $1
		 UNION
		 SELECT id, name
		 FROM expense_categories
		 WHERE user_id = $1`,
		sessionInfo.UserId,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var results []*CategoryWithName

	for rows.Next() {
		var c CategoryWithName
		if err := rows.Scan(
			&c.Id,
			&c.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, &c)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return results, nil
}

func (s *Service) GetCategoriesIdIncludes(
	ctx context.Context,
	ids []string,
) ([]*CategoryWithName, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
		args[i] = id
	}
	query := fmt.Sprintf(
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
		query,
		args...,
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close()

	var results []*CategoryWithName

	for rows.Next() {
		var c CategoryWithName
		if err := rows.Scan(
			&c.Id,
			&c.Name,
		); err != nil {
			return nil, err
		}
		results = append(results, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}
