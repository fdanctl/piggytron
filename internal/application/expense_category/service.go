package expensecategory

import (
	"context"
	"errors"
	"fmt"

	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/google/uuid"
)

var (
	ErrCategoryExists = errors.New("category already exists")
	ErrDuplicate      = errors.New("duplicate category")
)

type Service struct {
	repo expensecategory.Repository
}

func NewService(repo expensecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(
	ctx context.Context,
	name string,
	expenseType int,
) (*expensecategory.ExpenseCategory, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, errors.New("nil context")
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, errors.New("not sessionInfo")
	}

	et, err := expensecategory.NewExpenseType(uint8(expenseType))
	if err != nil {
		return nil, err
	}
	_, err = s.repo.FindByNameAndUser(ctx, expensecategory.ID(sessionInfo.UserId), name)
	if err == nil {
		return nil, ErrDuplicate
	}
	category, err := expensecategory.New(
		expensecategory.ID(uuid.New().String()),
		expensecategory.ID(sessionInfo.UserId),
		name,
		et,
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
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
		return nil, err
	}

	return categories, nil
}
