// Package appbudgethold implements the budgethold use cases
package appbudgethold

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/budgethold"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the budgethold use cases.
type Service struct {
	repo budgethold.Repository
}

// NewService wires the budget service to its repository.
func NewService(r budgethold.Repository) *Service {
	return &Service{repo: r}
}

// SaveBudgetHold creates a budget for a category and month, if not existing.
// Otherwise updates it.
func (s *Service) SaveBudgetHold(
	ctx context.Context,
	userID string,
	month budgethold.Month,
	amount int,
) (*budgethold.BudgetHold, error) {
	uid, err := util.ParseID[budgethold.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appbudgethold.SaveBudgetHold",
		)
		return nil, err
	}

	b, err := budgethold.New(
		uid,
		month,
		amount,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create budget hold",
			fmt.Errorf("failed to create budget: %w", err),
			"appbudgethold.SaveBudgetHold",
		)
		return nil, err
	}

	err = s.repo.Save(ctx, b)
	if err != nil {
		return nil, errs.NewInternalAppError(
			fmt.Errorf("failed saving budget hold: %w", err),
			"appbudgethold.SaveBudgetHold",
		)
	}
	return b, nil
}

// FindBudgetHold returns the budget hold for a month.
func (s *Service) FindBudgetHold(
	ctx context.Context,
	userID string,
	month budgethold.Month,
) (*budgethold.BudgetHold, error) {
	uid, err := util.ParseID[budgethold.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appbudget.FindBudget",
		)
		return nil, err
	}

	b, err := s.repo.FindByUserAndMonth(ctx, uid, month)
	if err != nil {
		if errors.Is(err, budgethold.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find budget hold",
				fmt.Errorf(
					"failed to find budget hold with user_id '%s' and month %s: %w",
					uid,
					month.String(),
					err,
				),
				"appbudgethold.FindBudgetHold",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf(
					"failed to find budget hold with user_id '%s' and month %s: %w",
					uid,
					month.String(),
					err,
				),
				"appbudgethold.FindBudgetHold",
			)
		}
		return nil, err
	}

	return b, nil
}
