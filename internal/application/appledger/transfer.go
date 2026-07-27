package appledger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

func (s *Service) CreateTransfer(
	ctx context.Context,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	categoryID string,
	srcAccID string,
	dstAccID string,
) (*ledger.Entry, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	id, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", srcAccID),
			fmt.Errorf("failed parsing id '%s': %w", srcAccID, err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	toAccID, err := util.ParseID[ledger.ID](dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", dstAccID),
			fmt.Errorf("failed parsing id '%s': %w", dstAccID, err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	var cid *ledger.ID
	if categoryID != "" {
		tempID, err := util.ParseID[ledger.ID](categoryID)
		if err != nil {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s is not a valid id", categoryID),
				fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
				"appledger.CreateTransfer",
			)
			return nil, err
		}
		cid = &tempID
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)

	fromAccount, err := qtx.GetAccountWithMinRunningBalance(ctx, srcAccID, date, nil, nil)
	if err != nil {
		return nil, errs.NewInternalAppError(
			fmt.Errorf("failed to find account '%s': %w", srcAccID, err),
			"appledger.CreateTransfer",
		)
	}
	toAccount, err := qtx.FindWithSum(ctx, dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Destination account not found",
			fmt.Errorf("failed to find account '%s': %w", dstAccID, err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	var accCID *ledger.ID
	if toAccount.Category.ID != util.ZeroUUID {
		temp := ledger.ID(toAccount.Category.ID)
		accCID = &temp
	}

	var toAccountCatType string
	if toAccount.IsSaving != nil && *toAccount.IsSaving && categoryID != "" {
		cattx := postgres.NewCategoryQueryService(tx)
		cat, err := cattx.FindByID(ctx, categoryID)
		if err != nil {
			errs.NewInternalAppError(err, "appledger.CreateTransfer")
			return nil, err
		}
		toAccountCatType = cat.Type
	}

	t, err := ledger.NewTransfer(
		id,
		uid,
		fromAccID,
		toAccID,
		cid,
		amount,
		description,
		date,
		fromAccount.MinRunningBalance,
		accCID,
		toAccountCatType,
		toAccount.IsSaving != nil && *toAccount.IsSaving,
	)
	if err != nil {
		msg := "Failed to create transfer"
		if errors.Is(err, ledger.ErrNegativeBalance) {
			msg = fmt.Sprintf(
				"%s becomes negative on %s",
				fromAccount.Name,
				fromAccount.MinDate.Format("January 2, 2006"),
			)
		}
		if errors.Is(err, ledger.ErrGoalCategory) {
			msg = fmt.Sprintf(
				"Transfers to %s must be have %s category",
				toAccount.Name,
				toAccount.Category.Name,
			)
		}
		if errors.Is(err, ledger.ErrNotSavingsCategory) {
			msg = "Category must be savings type to send money to savings account"
		}
		err = errs.NewAppError(
			errs.KindBusinessRule,
			msg,
			fmt.Errorf("failed to create income: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	// monthly summary
	// fromAcc
	fms, err := monthlysummary.New(
		monthlysummary.ID(fromAccID),
		monthlysummary.NewMonth(date),
		0,
		amount,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	err = mstx.Save(ctx, fms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	// toAcc
	tms, err := monthlysummary.New(
		monthlysummary.ID(toAccID),
		monthlysummary.NewMonth(date),
		amount,
		0,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	err = mstx.Save(ctx, tms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	return t, nil
}
