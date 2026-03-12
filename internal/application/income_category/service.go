package incomecategory

import (
	"context"
	"errors"
	"fmt"

	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

var ErrCategoryExists = errors.New("category already exists")

type Service struct {
	repo incomecategory.Repository
}

func NewService(repo incomecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, name string, expenseType int) error {
	fmt.Println(ctx.Value("user"))
	// _, err := s.repo.FindByNameAndUser()
	// if err == nil {
	// 	return err
	// }
	return nil
}

func (s *Service) ReadAllUserCategories(
	ctx context.Context,
) ([]*incomecategory.IncomeCategory, error) {
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

	categories, err := s.repo.FindAllByUser(ctx, incomecategory.ID(sessionInfo.UserId))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Println(categories)

	return categories, nil
}
