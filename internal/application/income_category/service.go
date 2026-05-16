package incomecategory

import (
	"context"

	incomecategory "github.com/fdanctl/piggytron/internal/domain/income_category"
	"github.com/google/uuid"
)

type Service struct {
	repo incomecategory.Repository
}

func NewService(repo incomecategory.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateCategory(
	ctx context.Context,
	userID string,
	name string,
) (*incomecategory.IncomeCategory, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	uid, err := incomecategory.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := incomecategory.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	category, err := incomecategory.New(id, uid, name)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		return nil, err
	}
	return category, nil
}

func (s *Service) ReadCategory(
	ctx context.Context,
	id string,
) (*incomecategory.IncomeCategory, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, incomecategory.ID(id))
}

func (s *Service) ReadAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*incomecategory.IncomeCategory, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	uid, err := incomecategory.NewID(userID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
