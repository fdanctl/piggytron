package appledger

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
	"github.com/fdanctl/piggytron/internal/util"
)

func (s *Service) CreateExpense(
	ctx context.Context,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	catID string,
	srcAccID string,
) (*ledger.Entry, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](catID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", cid),
			fmt.Errorf("failed parsing id '%s': %w", cid, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", fromAccID),
			fmt.Errorf("failed parsing id '%s': %w", fromAccID, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	id, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)

	acc, err := qtx.GetAccountWithMinRunningBalance(ctx, srcAccID, date, nil, nil)
	if err != nil {
		return nil, errs.NewInternalAppError(
			fmt.Errorf("failed to find account '%s': %w", srcAccID, err),
			"appledger.CreateExpense",
		)
	}

	var accCID *account.ID
	if acc.Category.ID != util.ZeroUUID {
		temp := account.ID(acc.Category.ID)
		accCID = &temp
	}

	a := account.Rehydrate(
		account.ID(acc.ID),
		account.ID(acc.UserID),
		account.AccountType(acc.Type),
		acc.Name,
		acc.IsSaving,
		acc.TargetAmount,
		acc.StartDate,
		acc.TargetDate,
		accCID,
		acc.Currency,
		acc.CreatedAt,
		acc.UpdatedAt,
	)

	if err := a.CanMakeExpense(); err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s of type %s can't make expenses", a.Name(), a.Type()),
			fmt.Errorf("%s of type %s can't make expenses: %w", a.Name(), a.Type(), err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	t, err := ledger.NewExpense(
		id,
		uid,
		fromAccID,
		cid,
		amount,
		description,
		date,
		acc.MinRunningBalance,
	)
	if err != nil {
		if errors.Is(err, ledger.ErrNegativeBalance) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf(
					"%s becomes negative on %s",
					acc.Name,
					acc.MinDate.Format("January 2, 2006"),
				),
				fmt.Errorf("failed to create expense: %w", err),
				"appledger.CreateExpense",
			)
		} else {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create expense",
				fmt.Errorf("failed to create expense: %w", err),
				"appledger.CreateExpense",
			)
		}
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	// monthly summary
	ms, err := monthlysummary.New(
		monthlysummary.ID(acc.ID),
		monthlysummary.NewMonth(date),
		0,
		amount,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	err = mstx.Save(ctx, ms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit: %w", err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	return t, nil
}
