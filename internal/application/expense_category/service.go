package expensecategory

import (
	"context"
	"errors"

	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
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
	userID string,
	name string,
	expenseType string,
) (*expensecategory.ExpenseCategory, error) {
	uid, err := expensecategory.NewID(userID)
	if err != nil {
		return nil, err
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindByNameAndUser(ctx, uid, name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := expensecategory.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	category, err := expensecategory.New(
		id,
		uid,
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
	userID string,
) ([]*expensecategory.ExpenseCategory, error) {
	uid, err := expensecategory.NewID(userID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
