package appledger

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
)

func (s *Service) updateMonthlySummary(
	ctx context.Context,
	mstx *postgres.MonthlySummaryRepository,
	inDelta, outDelta int,
	accId string,
	month monthlysummary.Month,
) error {
	// monthly summary — read current, add/subtract, write back
	fms, err := mstx.FindByAccountAndMonth(ctx, accId, month)
	if err != nil {
		if errors.Is(err, monthlysummary.ErrNotFound) {
			fms, err := monthlysummary.New(
				monthlysummary.ID(accId),
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
