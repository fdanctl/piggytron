package appaccount

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/util"
)

func (s *Service) CompleteGoal(
	ctx context.Context,
	aid string,
	userID string,
	amount int,
	date time.Time,
	remainingAccID string,
) (*account.Account, error) {
	logger := middleware.LoggerFromContext(ctx)
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](aid)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", aid),
			fmt.Errorf("failed parsing id '%s': %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	goalAcc, err := qtx.FindWithSum(ctx, aid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding ledger entries for '%s' account: %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	if goalAcc.UserID != userID {
		return nil, account.ErrNotFound
	}

	logger.Debug(
		"balance",
		"balance",
		goalAcc.Sum,
		"amount",
		amount,
		"remaining",
		goalAcc.Sum-amount,
	)

	if goalAcc.Sum-amount < 0 {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Amount used can't be higher than balance",
			fmt.Errorf("amount used can't be higher than balance: %w", account.ErrNegativeBalance),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	lastEntryDate, err := qtx.GetAccountLastEntryDate(ctx, aid)
	if date.Before(lastEntryDate) {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Contributions found after the complete date",
			fmt.Errorf("contributions found after the complete date: %w", account.ErrInvalidDate),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	// remaining transfer
	if goalAcc.Sum-amount > 0 {
		if remainingAccID == "" {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination account is necessary for the remaining",
				fmt.Errorf(
					"empty account id: %w",
					account.ErrInvalidID,
				),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}
		remainingTarget, err := qtx.FindWithSum(ctx, remainingAccID)
		if err != nil {
			if errors.Is(err, account.ErrNotFound) {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					"Remaining destination account not found",
					fmt.Errorf(
						"remaining destination account not found: %w",
						err,
					),
					"appaccount.CompleteGoal",
				)
			} else {
				err = errs.NewInternalAppError(err, "appaccount.CompleteGoal")
			}
			return nil, err
		}

		if remainingTarget.Type != string(account.CheckingType) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining amount can only be transfer to a checking account",
				fmt.Errorf(
					"remaining amount can only be transfer to a checking account: %w",
					account.ErrNegativeBalance,
				),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		if remainingTarget.Status == string(account.ClosedStatus) {
			return nil, errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination is closed",
				account.ErrClosedAccount,
				"appaccount.CompleteGoal",
			)
		}

		tid, err := util.NewID[ledger.ID]()
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed generating id: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		toAccID, err := util.ParseID[ledger.ID](remainingTarget.ID)
		if err != nil {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s is not a valid id", remainingTarget.ID),
				fmt.Errorf("failed parsing id '%s': %w", remainingTarget.ID, err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		t, err := ledger.NewTransfer(
			tid,
			uid,
			fromAccID,
			toAccID,
			nil,
			goalAcc.Sum-amount,
			fmt.Sprintf("Remaining from completing %s goal", goalAcc.Name),
			date,
			goalAcc.Sum,
			nil,
			"",
			remainingTarget.Type == string(account.SavingsType),
		)
		if err != nil {
			msg := "Failed to create transfer"
			if errors.Is(err, ledger.ErrNegativeBalance) {
				msg = fmt.Sprintf(
					"%s becomes negative on %s",
					goalAcc.Name,
					date,
				)
			}
			err = errs.NewAppError(
				errs.KindBusinessRule,
				msg,
				fmt.Errorf("failed to create transfer: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		err = rtx.Create(ctx, t)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving transaction: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		// monthly summary
		// toAcc
		tms, err := monthlysummary.New(
			monthlysummary.ID(remainingTarget.ID),
			monthlysummary.NewMonth(date),
			goalAcc.Sum-amount,
			0,
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create summary",
				fmt.Errorf("failed to create summary: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		err = mstx.Save(ctx, tms)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving summary: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}
	}

	cid := &goalAcc.Category.ID
	if *cid == util.ZeroUUID {
		cid = nil
	}

	acc := account.Rehydrate(
		account.ID(goalAcc.ID),
		account.ID(goalAcc.UserID),
		account.AccountType(goalAcc.Type),
		goalAcc.Name,
		account.AccountStatus(goalAcc.Status),
		goalAcc.TargetAmount,
		goalAcc.StartDate,
		goalAcc.TargetDate,
		(*account.ID)(cid),
		goalAcc.CompletedAt,
		goalAcc.CancelledAt,
		goalAcc.Currency,
		goalAcc.ClosedAt,
		goalAcc.CreatedAt,
		goalAcc.UpdatedAt,
	)

	err = acc.CompleteGoal()
	if err != nil {
		return nil, err
	}

	atx := postgres.NewAccountRepository(tx)
	err = atx.UpdateStatus(ctx, acc)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to close '%s' account: %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	gfid, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	fulfillmentEntry, err := ledger.NewGoalFulfillment(
		gfid,
		uid,
		fromAccID,
		acc.Name(),
		ledger.ID(*acc.CategoryID()),
		amount,
		date,
	)
	logger.Debug("fulfillment entry", "entry", fulfillmentEntry)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create transfer",
			fmt.Errorf("failed to create goal fulfillment: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	err = rtx.Create(ctx, fulfillmentEntry)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	// monthly summary
	// fromAcc
	fms, err := monthlysummary.New(
		monthlysummary.ID(aid),
		monthlysummary.NewMonth(date),
		0,
		goalAcc.Sum,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	err = mstx.Save(ctx, fms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	return acc, nil
}
