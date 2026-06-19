package appexpensecategory

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/errs"
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
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	et, err := expensecategory.NewExpenseType(expenseType)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s is not a valid expense type", et),
			fmt.Errorf("%s is not a valid expense type: %w", et, err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	id, err := util.NewID[expensecategory.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	category, err := expensecategory.New(id, uid, name, et)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create category",
			fmt.Errorf("failed to create category: %w", err),
			"appexpensecategory.CreateCategory",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, category)
	if err != nil {
		if errors.Is(err, expensecategory.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindValidation,
				"An expense category with the same name already exists",
				fmt.Errorf("failed saving category '%s': %w", category.Name(), err),
				"appexpensecategory.CreateCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving category: %w", err),
				"appexpensecategory.CreateCategory",
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
) (*expensecategory.ExpenseCategory, error) {
	cid, err := util.ParseID[expensecategory.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}

	cat, err := s.repo.FindByID(ctx, cid)
	if err != nil {
		if errors.Is(err, expensecategory.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The category does not exists",
				fmt.Errorf("failed to found category '%s': %w", id, err),
				"appexpensecategory.FindCategory",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", id, err),
				"appexpensecategory.FindCategory",
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
				expensecategory.ErrNotFound,
			),
			"appexpensecategory.FindCategory",
		)
		return nil, err
	}

	return cat, nil
}

func (s *Service) FindAllUserCategories(
	ctx context.Context,
	userID string,
) ([]*expensecategory.ExpenseCategory, error) {
	uid, err := util.ParseID[expensecategory.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appexpensecategory.FindAllUserCategories",
		)
		return nil, err
	}

	categories, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' expense categories: %w", uid, err),
			"appexpensecategory.FindAllUserCategories",
		)
		return nil, err
	}

	return categories, nil
}
