package appledger

import (
	"cmp"
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

func (s *Service) CreateInterest(
	ctx context.Context,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	categoryID string,
	dstAccID string,
) (*ledger.Entry, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	id, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	toAccID, err := util.ParseID[ledger.ID](dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", dstAccID),
			fmt.Errorf("failed parsing id '%s': %w", dstAccID, err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}
	defer tx.Rollback()

	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	ictx := postgres.NewIncomeCategoryRepository(tx)

	_, err = isIncomeCategoryValid(ctx, ictx, categoryID, date)
	if err != nil {
		return nil, err
	}

	t, err := ledger.NewInterest(
		id,
		uid,
		toAccID,
		cid,
		amount,
		description,
		date,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create interest",
			fmt.Errorf("failed to create interest: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	// monthly summary
	ms, err := monthlysummary.New(
		monthlysummary.ID(dstAccID),
		monthlysummary.NewMonth(date),
		amount,
		0,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	err = mstx.Save(ctx, ms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit: %w", err),
			"appledger.CreateInterest",
		)
		return nil, err
	}

	return t, nil
}

func (s *Service) UpdateInterest(
	ctx context.Context,
	id string,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	categoryID string,
	dstAccID string,
) (*ledger.Entry, error) {
	tid, err := util.ParseID[ledger.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	toAccID, err := util.ParseID[ledger.ID](dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", dstAccID),
			fmt.Errorf("failed parsing id '%s': %w", dstAccID, err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	ictx := postgres.NewIncomeCategoryRepository(tx)

	_, err = isIncomeCategoryValid(ctx, ictx, categoryID, date)
	if err != nil {
		return nil, err
	}

	t, err := rtx.FindByID(ctx, tid)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find transaction",
				fmt.Errorf("failed to find transaction '%s': %w", tid, err),
				"appledger.UpdateInterest",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find transaction '%s': %w", id, err),
				"appledger.UpdateInterest",
			)
		}
		return nil, err
	}

	if t.UserID() != uid || t.Type() != "interest" {
		fmt.Println(t.Type())
		return nil, ledger.ErrNotFound
	}

	accountChanged := t.ToAccountID() == nil || string(*t.ToAccountID()) != string(toAccID)
	prevAmount := t.Amount()
	prevDate := t.Date()
	prevAccountID := string(*t.ToAccountID())
	if err := t.UpdateInterest(toAccID, cid, amount, description, date); err != nil {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to update interest",
			fmt.Errorf("failed to update interest: %w", err),
			"appledger.UpdateInterest",
		)
	}

	// verify if the account is closed
	oldAcc, err := qtx.FindWithSum(ctx, prevAccountID)
	if err != nil {
		return nil, errs.NewInternalAppError(err, "appledger.UpdateInterest")
	}
	var oldAccCID *account.ID
	if oldAcc.Category.ID != util.ZeroUUID {
		temp := account.ID(oldAcc.Category.ID)
		oldAccCID = &temp
	}
	oldA := account.Rehydrate(
		account.ID(oldAcc.ID),
		account.ID(oldAcc.UserID),
		account.AccountType(oldAcc.Type),
		oldAcc.Name,
		account.AccountStatus(oldAcc.Status),
		oldAcc.TargetAmount,
		oldAcc.StartDate,
		oldAcc.TargetDate,
		oldAccCID,
		oldAcc.CompletedAt,
		oldAcc.CancelledAt,
		oldAcc.FinalizedAmount,
		oldAcc.Currency,
		oldAcc.ClosedAt,
		oldAcc.CreatedAt,
		oldAcc.UpdatedAt,
	)

	if oldA.IsClosed() {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Can't change closed accounts",
			account.ErrClosedAccount,
			"appledger.UpdateInterest",
		)
	}

	// Account changed: PREV account loses income → verify PREV account stays solvent.
	if accountChanged {
		// Verify the PREV account (it loses this income)
		acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevAccountID, prevDate, nil, &id)
		if err != nil {
			return nil, errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
				"appledger.UpdateInterest",
			)
		}

		if acc.MinRunningBalance < 0 {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf(
					"%s becomes negative on %s",
					acc.Name,
					acc.MinDate.Format("January 2, 2006"),
				),
				fmt.Errorf("balance check failed: %w", ledger.ErrNegativeBalance),
				"appledger.UpdateInterest",
			)
			return nil, err
		}
	} else {
		// Same account
		c, err := getCaseFromTable(date.Compare(prevDate), cmp.Compare(amount, prevAmount))
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("don't know how could I make this error: %w", err),
				"appledger.UpdateInterest",
			)
			return nil, err
		}
		// 1. Lower date && lower amount
		// 2. Lower date && same amount
		// 3. Lower date && higher amount
		// 4. Same date && lower amount
		// 5. Same date && same amount
		// 6. Same date && higher amount
		// 7. Higher date && lower amount
		// 8. Higher date && same amount
		// 9. Higher date && higher amount
		switch c {
		// Same account — only verify if the update is "worse" than before:
		// not less money and money does not arrives later.
		case 2, 3, 5, 6:
			// SAFE
		case 1, 4:
			// Check from PREV date the balance without the PREV amount + the NEW amount
			acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevAccountID, prevDate, nil, &id)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
					"appledger.UpdateInterest",
				)
			}
			if acc.MinRunningBalance+amount < 0 {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					fmt.Sprintf(
						"%s becomes negative on %s",
						acc.Name,
						acc.MinDate.Format("January 2, 2006"),
					),
					fmt.Errorf("case %d, balance check failed: %w", c, ledger.ErrNegativeBalance),
					"appledger.UpdateInterest",
				)
				return nil, err
			}

		case 7, 8, 9:
			// Check from PREV date to NEW date (not included) the balance without the PREV amount
			until := time.Date(date.Year(), date.Month(), date.Day()-1, 0, 0, 0, 0, date.Location())
			acc, err := qtx.GetAccountWithMinRunningBalance(
				ctx,
				prevAccountID,
				prevDate,
				&until,
				&id,
			)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
					"appledger.UpdateInterest",
				)
			}
			if acc.MinRunningBalance < 0 {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					fmt.Sprintf(
						"%s becomes negative on %s",
						acc.Name,
						acc.MinDate.Format("January 2, 2006"),
					),
					fmt.Errorf("case %d, balance check failed: %w", c, ledger.ErrNegativeBalance),
					"appledger.UpdateInterest",
				)
				return nil, err
			}
			if c == 7 {
				// Check from NEW date the balance without the PREV amount + the NEW amount
				acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevAccountID, date, nil, &id)
				if err != nil {
					return nil, errs.NewInternalAppError(
						fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
						"appledger.UpdateInterest",
					)
				}
				if acc.MinRunningBalance+amount < 0 {
					err = errs.NewAppError(
						errs.KindBusinessRule,
						fmt.Sprintf(
							"%s becomes negative on %s",
							acc.Name,
							acc.MinDate.Format("January 2, 2006"),
						),
						fmt.Errorf(
							"case %d, balance check failed: %w",
							c,
							ledger.ErrNegativeBalance,
						),
						"appledger.UpdateInterest",
					)
					return nil, err
				}
			}
		}
	}

	if err := rtx.Update(ctx, t); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	// Monthly summary bookkeeping.
	newAccountID := string(toAccID)

	if accountChanged {
		// Case A: account changed.
		// PREV account: remove PREV income from PREV month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			-prevAmount,
			0,
			prevAccountID,
			monthlysummary.NewMonth(prevDate),
		); err != nil {
			return nil, err
		}
		// NEW account: add NEW income at NEW month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			amount,
			0,
			newAccountID,
			monthlysummary.NewMonth(date),
		); err != nil {
			return nil, err
		}
	} else {
		// Case B: same account.
		prevMonth := monthlysummary.NewMonth(prevDate)
		newMonth := monthlysummary.NewMonth(date)

		if prevMonth.Time().Compare(newMonth.Time()) != 0 {
			// Different months: remove PREV income from PREV month, add NEW income to NEW month.
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				-prevAmount,
				0,
				newAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				amount,
				0,
				newAccountID,
				newMonth,
			); err != nil {
				return nil, err
			}
		} else {
			// Same month: net delta = newAmount - prevAmount.
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				amount-prevAmount,
				0,
				newAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit: %w", err),
			"appledger.UpdateInterest",
		)
		return nil, err
	}

	return t, nil
}
