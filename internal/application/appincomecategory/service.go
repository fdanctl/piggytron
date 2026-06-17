package appincomecategory

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/fdanctl/piggytron/internal/util"
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
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		return nil, err
	}

	id, err := util.NewID[incomecategory.ID]()
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

func (s *Service) FindCategory(
	ctx context.Context,
	id string,
) (*incomecategory.IncomeCategory, error) {
	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, cid)
}

func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*incomecategory.IncomeCategory, error) {
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
