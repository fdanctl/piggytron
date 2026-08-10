// Package appbudget implements the budget use cases: creating and querying
// monthly budgets per expense category, and copying the previous month's
// budgets into an empty month.
package appbudget

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the budget use cases.
type Service struct {
	repo budget.Repository
}

// NewService wires the budget service to its repository.
func NewService(r budget.Repository) *Service {
	return &Service{repo: r}
}

// SaveBudget creates a budget for a category and month, if not existing.
// Otherwise updates it.
func (s *Service) SaveBudget(
	ctx context.Context,
	categoryID string,
	month budget.Month,
	amount int,
) (*budget.Budget, error) {
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

	b, err := budget.New(
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

	err = s.repo.Save(ctx, b)
	if err != nil {
		return nil, errs.NewInternalAppError(
			fmt.Errorf("failed saving budget: %w", err),
			"appbudget.CreateBudget",
		)
	}
	return b, nil
}

// FindBudget returns the budget for a category and month.
func (s *Service) FindBudget(
	ctx context.Context,
	categoryID string,
	month budget.Month,
) (*budget.Budget, error) {
	cid, err := util.ParseID[budget.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appbudget.FindBudget",
		)
		return nil, err
	}

	b, err := s.repo.FindByCategoryAndMonth(ctx, cid, month)
	if err != nil {
		if errors.Is(err, budget.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find budget",
				fmt.Errorf(
					"failed to find budget with id '%s' and month %s: %w",
					cid,
					month.String(),
					err,
				),
				"appbudget.FindBudget",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf(
					"failed to find budget with id '%s' and month %s: %w",
					cid,
					month.String(),
					err,
				),
				"appbudget.FindBudget",
			)
		}
		return nil, err
	}

	return b, nil
}

// CopyFromLastMonth carries the previous month's budgets into month,
// returning the number of categories updated.
func (s *Service) CopyFromLastMonth(
	ctx context.Context,
	userID string,
	month budget.Month,
) (int, error) {
	uid, err := util.ParseID[budget.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appbudget.CopyFromLastMonth",
		)
		return 0, err
	}

	count, err := s.repo.CopyLastMonthBudget(ctx, uid, month)
	if err != nil {
		return 0, errs.NewInternalAppError(err, "appbudget.CopyFromLastMonth")
	}
	return count, nil
}
