package appbudget

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/util"
)

type Service struct {
	repo budget.Repository
}

func NewService(r budget.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateBudget(
	ctx context.Context,
	userID string,
	categoryID string,
	month time.Time,
	amount int,
) (*budget.Budget, error) {
	uid, err := util.ParseID[budget.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appbudget.CreateBudget",
		)
		return nil, err
	}

	cid, err := util.ParseID[budget.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", cid),
			fmt.Errorf("failed parsing id '%s': %w", cid, err),
			"appbudget.CreateBudget",
		)
		return nil, err
	}

	id, err := util.NewID[budget.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appbudget.CreateBudget",
		)
		return nil, err
	}

	b, err := budget.New(
		id,
		uid,
		cid,
		month,
		amount,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create budget",
			fmt.Errorf("failed to create budget: %w", err),
			"appbudget.CreateBudget",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, b)
	if err != nil {
		if errors.Is(err, budget.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf("Already exist a budget for %s", month.String()),
				fmt.Errorf("already exist a %s budget for %s", cid, month.String()),
				"appbudget.CreateBudget",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving budget: %w", err),
				"appbudget.CreateBudget",
			)
		}
		return nil, err
	}
	return b, nil
}

func (s *Service) UpdateBudgetAmount(
	ctx context.Context,
	id string,
	amount int,
) error {
	bid, err := util.ParseID[budget.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", bid),
			fmt.Errorf("failed parsing id '%s': %w", bid, err),
			"appbudget.UpdateBudgetAmount",
		)
		return err
	}

	if amount < 0 {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Amount can't be negative",
			fmt.Errorf("%d is not valid: %w", amount, budget.ErrInvalidAmount),
			"appbudget.UpdateBudgetAmount",
		)
		return err
	}

	err = s.repo.UpdateAmount(ctx, bid, amount)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to update amount for '%s' budget: %w", bid, err),
			"appbudget.UpdateBudgetAmount",
		)
		return err
	}
	return nil
}

func (s *Service) FindBudget(
	ctx context.Context,
	id string,
) (*budget.Budget, error) {
	bid, err := util.ParseID[budget.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", bid),
			fmt.Errorf("failed parsing id '%s': %w", bid, err),
			"appbudget.FindBudget",
		)
		return nil, err
	}

	b, err := s.repo.FindByID(ctx, bid)
	if err != nil {
		if errors.Is(err, budget.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find budget",
				fmt.Errorf("failed to find budget '%s': %w", bid, err),
				"appbudget.FindBudget",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find budget '%s': %w", bid, err),
				"appbudget.FindBudget",
			)
		}
		return nil, err
	}

	return b, nil
}
