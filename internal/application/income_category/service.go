package incomecategory

import (
	"context"
	"errors"
	"fmt"

	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/google/uuid"
)

var (
	ErrCategoryExists = errors.New("category already exists")
	ErrDuplicate      = errors.New("duplicate category")
)

type Service struct {
	repo incomecategory.Repository
}

func NewService(repo incomecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(ctx context.Context, name string) error {
	v := ctx.Value("user")
	if v == nil {
		fmt.Println("nil context")
		return nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		fmt.Println("not sessionInfo")
		return nil
	}

	_, err := s.repo.FindByNameAndUser(ctx, incomecategory.ID(sessionInfo.UserId), name)
	if err == nil {
		return ErrDuplicate
	}
	u, err := incomecategory.New(
		incomecategory.ID(uuid.New().String()),
		incomecategory.ID(sessionInfo.UserId),
		name,
	)
	if err != nil {
		return err
	}

	err = s.repo.Save(ctx, u)
	if err != nil {
		return err
	}
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

	return categories, nil
}
