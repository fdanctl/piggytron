package appexpensecategory

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/util"
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
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		return nil, err
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		return nil, err
	}

	id, err := util.NewID[expensecategory.ID]()
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

	err = s.repo.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) FindCategory(
	ctx context.Context,
	id string,
) (*expensecategory.ExpenseCategory, error) {
	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, cid)
}

func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*expensecategory.ExpenseCategory, error) {
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
