package expensecategory

import (
	"context"
	"errors"
	"fmt"

	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
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
	expenseType string,
) (*expensecategory.ExpenseCategory, error) {
	v := ctx.Value(middleware.UserKey)
	if v == nil {
		return nil, errors.New("nil context")
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, errors.New("not sessionInfo")
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		return nil, err
	}
	_, err = s.repo.FindByNameAndUser(ctx, expensecategory.ID(sessionInfo.UserID), name)
	if err == nil {
		return nil, ErrDuplicate
	}
	category, err := expensecategory.New(
		expensecategory.ID(uuid.New().String()),
		expensecategory.ID(sessionInfo.UserID),
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

func (s *Service) ReadCategory(
	ctx context.Context,
	id string,
) (*expensecategory.ExpenseCategory, error) {
	return s.repo.FindByID(ctx, expensecategory.ID(id))
}

func (s *Service) ReadAllUserCategories(
	ctx context.Context,
) ([]*expensecategory.ExpenseCategory, error) {
	v := ctx.Value(middleware.UserKey)
	if v == nil {
		fmt.Println("nil context")
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		fmt.Println("not sessionInfo")
		return nil, nil
	}

	categories, err := s.repo.FindAllByUser(ctx, expensecategory.ID(sessionInfo.UserID))
	if err != nil {
		return nil, err
	}

	return categories, nil
}
