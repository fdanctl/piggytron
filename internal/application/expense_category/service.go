package expensecategory

import (
	"context"
	"errors"
	"fmt"

	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

var ErrCategoryExists = errors.New("category already exists")

type Service struct {
	repo expensecategory.Repository
}

func NewService(repo expensecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, name string, expenseType int) error {
	v := ctx.Value("user")
	if v == nil {
		return nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil
	}
	fmt.Println(sessionInfo.UserId)
	return nil
}

func (s *Service) ReadAllUserCategories(
	ctx context.Context,
) ([]*expensecategory.ExpenseCategory, error) {
	v := ctx.Value("user")
	if v == nil {
		fmt.Println("nil context")
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		fmt.Println("not sessionInfo")
		return nil, nil
	}

	categories, err := s.repo.FindAllByUser(ctx, expensecategory.ID(sessionInfo.UserId))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(categories)

	return categories, nil
}
