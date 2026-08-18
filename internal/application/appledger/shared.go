package appledger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/expensecategory"
	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

func (s *Service) updateMonthlySummary(
	ctx context.Context,
	mstx *postgres.MonthlySummaryRepository,
	inDelta, outDelta int,
	accID string,
	month monthlysummary.Month,
) error {
	// monthly summary — read current, add/subtract, write back
	fms, err := mstx.FindByAccountAndMonth(ctx, accID, month)
	if err != nil {
		if errors.Is(err, monthlysummary.ErrNotFound) {
			fms, err := monthlysummary.New(
				monthlysummary.ID(accID),
				month,
				inDelta,
				outDelta,
			)
			if err != nil {
				err = errs.NewAppError(
					errs.KindInternal,
					"Failed to create summary",
					fmt.Errorf("failed to create summary: %w", err),
					"appledger.updateMonthlySummary",
				)
				return err
			}

			err = mstx.Save(ctx, fms)
			if err != nil {
				err = errs.NewInternalAppError(
					fmt.Errorf("failed saving summary: %w", err),
					"appledger.updateMonthlySummary",
				)
				return err
			}
			return nil
		}

		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding summary: %w", err),
			"appledger.updateMonthlySummary",
		)
		return err
	}
	if err := fms.AddMoneyIn(inDelta); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating summary: %w", err),
			"appledger.updateMonthlySummary",
		)
		return err
	}
	if err := fms.AddMoneyOut(outDelta); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating summary: %w", err),
			"appledger.updateMonthlySummary",
		)
		return err
	}
	if err := mstx.Update(ctx, fms); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appledger.updateMonthlySummary",
		)
		return err
	}
	return nil
}

// getCaseFromTable return 1 out of 9 cases from a possiblity table (cond1 × cond2)
//
//	| ↓cond1 x cond2→ | -   | =   | +   |
//	| -               | 1   | 2   | 3   |
//	| =               | 4   | 5   | 6   |
//	| +               | 7   | 8   | 9   |
//
// cond must be -1 for negative, 0 for equal, and 1 for positive
func getCaseFromTable(cond1, cond2 int) (int, error) {
	if cond1 < -1 || cond1 > 1 {
		return 0, fmt.Errorf("%d is invalid: %w", cond1, errors.New("invalid conditional result"))
	}
	if cond2 < -1 || cond2 > 1 {
		return 0, fmt.Errorf("%d is invalid: %w", cond2, errors.New("invalid conditional result"))
	}
	return (3 * cond1) + cond2 + 5, nil
}

func isIncomeCategoryValid(
	ctx context.Context,
	tx incomecategory.Repository,
	categoryID string,
	date time.Time,
) (bool, error) {
	cid, err := util.ParseID[incomecategory.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appledger.isIncomeCategoryValid",
		)
		return false, err
	}

	cat, err := tx.FindByID(ctx, cid)
	if err != nil {
		err = errs.NewInternalAppError(
			err,
			"appledger.isIncomeCategoryValid",
		)
	}
	if cat.Status() != incomecategory.ActiveStatus && date.After(*cat.ArchivedAt()) {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf(
				"%s category was archived before %s",
				cat.Name(),
				date.Format(time.DateOnly),
			),
			errors.New("category was archived before entry data"),
			"appledger.isIncomeCategoryValid",
		)
		return false, err
	}

	return true, nil
}

func isExpenseCategoryValid(
	ctx context.Context,
	tx expensecategory.Repository,
	categoryID string,
	date time.Time,
) (bool, error) {
	cid, err := util.ParseID[expensecategory.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appledger.isExpenseCategoryValid",
		)
		return false, err
	}

	cat, err := tx.FindByID(ctx, cid)
	if err != nil {
		err = errs.NewInternalAppError(
			err,
			"appledger.isIncomeCategoryValid",
		)
	}
	if cat.Status() != expensecategory.ActiveStatus && date.After(*cat.ArchivedAt()) {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf(
				"%s category was archived before %s",
				cat.Name(),
				date.Format(time.DateOnly),
			),
			errors.New("category was archived before entry data"),
			"appledger.isExpenseCategoryValid",
		)
		return false, err
	}

	return true, nil
}
