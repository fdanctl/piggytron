package appincomecategory

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/fdanctl/piggytron/internal/errs"
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
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	id, err := util.NewID[incomecategory.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	category, err := incomecategory.New(id, uid, name)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appincomecategory.CreateCategory",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, incomecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindValidation,
				"An income category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appincomecategory.CreateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appincomecategory.CreateCategory",
			)
		}
		return nil, err
	}
	return category, nil
}

func (s *Service) FindCategory(
	ctx context.Context,
	id string,
	userID string,
) (*incomecategory.IncomeCategory, error) {
	cid, err := util.ParseID[incomecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}

	cat, err := s.repo.FindByID(ctx, cid)
	if err != nil {
		if errors.Is(err, incomecategory.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The category does not exists",
				fmt.Errorf("failed to found category '%s': %w", id, err),
				"appincomecategory.FindCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", id, err),
				"appincomecategory.FindCategory",
			)
		}
		return nil, err
	}

	if cat.UserID() != uid {
		err = errs.NewAppError(
			errs.KindNotFound,
			"The category does not exists",
			fmt.Errorf(
				"the category does not belong to user '%s': %w",
				uid,
				incomecategory.ErrNotFound,
			),
			"appincomecategory.FindCategory",
		)
		return nil, err
	}

	return cat, nil
}

func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*incomecategory.IncomeCategory, error) {
	uid, err := util.ParseID[incomecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appincomecategory.FindAllUserCategories",
		)
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' income categories: %w", uid, err),
			"appincomecategory.FindAllUserCategories",
		)
		return nil, err
	}

	return categories, nil
}
