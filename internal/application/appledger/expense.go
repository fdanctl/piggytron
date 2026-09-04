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

// CreateExpense records money spent from srcAccID, refusing goals and savings
// accounts, rejecting overdraft, and adds the amount to the monthly summary.
func (s *Service) CreateExpense(
	ctx context.Context,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	note string,
	catID string,
	srcAccID string,
) (*ledger.Entry, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](catID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", catID),
			fmt.Errorf("failed parsing id '%s': %w", catID, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", srcAccID),
			fmt.Errorf("failed parsing id '%s': %w", srcAccID, err),
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
	ltx := postgres.NewLedgerQueryService(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	ectx := postgres.NewExpenseCategoryRepository(tx)

	_, err = isExpenseCategoryValid(ctx, ectx, catID, date)
	if err != nil {
		return nil, err
	}

	err = s.checkAndUpdateInitialBalance(ctx, userID, rtx, ltx, mstx, srcAccID, date)
	if err != nil {
		return nil, err
	}

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
		account.AccountStatus(acc.Status),
		acc.TargetAmount,
		acc.StartDate,
		acc.TargetDate,
		accCID,
		acc.Note,
		acc.CompletedAt,
		acc.CancelledAt,
		acc.FinalizedAmount,
		acc.Currency,
		acc.ClosedAt,
		acc.CreatedAt,
		acc.UpdatedAt,
	)

	if err := a.CanMakeExpense(); err != nil {
		msg := fmt.Sprintf("%s of type %s can't make expenses", a.Name(), a.Type())
		if errors.Is(err, account.ErrClosedAccount) {
			msg = fmt.Sprintf("%s is closed", a.Name())
		}
		err = errs.NewAppError(
			errs.KindBusinessRule,
			msg,
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
		note,
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

// UpdateExpense modifies an expense entry, re-running the overdraft checks
// affected by the change and reconciling monthly summaries.
func (s *Service) UpdateExpense(
	ctx context.Context,
	id string,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	note string,
	categoryID string,
	srcAccID string,
) (*ledger.Entry, error) {
	tid, err := util.ParseID[ledger.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", srcAccID),
			fmt.Errorf("failed parsing id '%s': %w", srcAccID, err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	ltx := postgres.NewLedgerQueryService(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	ectx := postgres.NewExpenseCategoryRepository(tx)

	_, err = isExpenseCategoryValid(ctx, ectx, categoryID, date)
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
				"appledger.UpdateExpense",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find transaction '%s': %w", id, err),
				"appledger.UpdateExpense",
			)
		}
		return nil, err
	}

	if t.UserID() != uid || t.Type() != "expense" {
		return nil, ledger.ErrNotFound
	}

	accountChanged := t.FromAccountID() == nil || string(*t.FromAccountID()) != string(fromAccID)
	prevAmount := t.Amount()
	prevDate := t.Date()
	prevAccountID := string(*t.FromAccountID())

	if accountChanged {
		err = s.checkAndUpdateInitialBalance(ctx, userID, rtx, ltx, mstx, srcAccID, date)
		if err != nil {
			return nil, err
		}
	} else {
		err = s.checkAndUpdateInitialBalance(ctx, userID, rtx, ltx, mstx, prevAccountID, date)
		if err != nil {
			return nil, err
		}
	}

	if err := t.UpdateExpense(fromAccID, cid, amount, description, date, note); err != nil {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to update expense",
			fmt.Errorf("failed to update expense: %w", err),
			"appledger.UpdateExpense",
		)
	}

	// verify if the account is closed
	oldAcc, err := qtx.FindWithSum(ctx, prevAccountID)
	if err != nil {
		return nil, errs.NewInternalAppError(err, "appledger.UpdateExpense")
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
		oldAcc.Note,
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
			"appledger.UpdateExpense",
		)
	}

	// Account changed: NEW account gains expense → verify NEW account stays solvent.
	if accountChanged {
		// Verify the NEW account (it gains expense)
		acc, err := qtx.GetAccountWithMinRunningBalance(ctx, srcAccID, date, nil, nil)
		if err != nil {
			return nil, errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
				"appledger.UpdateExpense",
			)
		}

		//////
		// check if new account can make expense
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
			account.AccountStatus(acc.Status),
			acc.TargetAmount,
			acc.StartDate,
			acc.TargetDate,
			accCID,
			acc.Note,
			acc.CompletedAt,
			acc.CancelledAt,
			acc.FinalizedAmount,
			acc.Currency,
			acc.ClosedAt,
			acc.CreatedAt,
			acc.UpdatedAt,
		)

		if err := a.CanMakeExpense(); err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf("%s of type %s can't make expenses", a.Name(), a.Type()),
				fmt.Errorf("%s of type %s can't make expenses: %w", a.Name(), a.Type(), err),
				"appledger.UpdateExpense",
			)
			return nil, err
		}
		//////

		if acc.MinRunningBalance-amount < 0 {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf(
					"%s becomes negative on %s",
					acc.Name,
					acc.MinDate.Format("January 2, 2006"),
				),
				fmt.Errorf("balance check failed: %w", ledger.ErrNegativeBalance),
				"appledger.UpdateExpense",
			)
			return nil, err
		}
	} else {
		// Same account
		c, err := getCaseFromTable(date.Compare(prevDate), cmp.Compare(amount, prevAmount))
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("don't know how could I make this error: %w", err),
				"appledger.UpdateExpense",
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
		// not more money and money does not arrives before.
		case 4, 5, 7, 8:
			// SAFE
		case 1, 2, 3:
			// Check from NEW date until PREV date the balance - NEW amount >= 0
			until := time.Date(
				prevDate.Year(),
				prevDate.Month(),
				prevDate.Day()-1,
				0,
				0,
				0,
				0,
				prevDate.Location(),
			)
			acc, err := qtx.GetAccountWithMinRunningBalance(ctx, srcAccID, date, &until, nil)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
					"appledger.UpdateExpense",
				)
			}
			if acc.MinRunningBalance-amount < 0 {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					fmt.Sprintf(
						"%s becomes negative on %s",
						acc.Name,
						acc.MinDate.Format("January 2, 2006"),
					),
					fmt.Errorf("case %d, balance check failed: %w", c, ledger.ErrNegativeBalance),
					"appledger.UpdateExpense",
				)
				return nil, err
			}
			if c == 3 {
				// Check from PREV date the balance + PREV amount - NEW amount >= 0
				acc, err := qtx.GetAccountWithMinRunningBalance(
					ctx,
					prevAccountID,
					prevDate,
					nil,
					nil,
				)
				if err != nil {
					return nil, errs.NewInternalAppError(
						fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
						"appledger.UpdateExpense",
					)
				}
				if acc.MinRunningBalance+prevAmount-amount < 0 {
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
						"appledger.UpdateExpense",
					)
					return nil, err
				}
			}

		case 6, 9:
			// Check from NEW date the balance + PREV amount - NEW amount >= 0
			acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevAccountID, date, nil, nil)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevAccountID, err),
					"appledger.UpdateExpense",
				)
			}
			if acc.MinRunningBalance+prevAmount-amount < 0 {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					fmt.Sprintf(
						"%s becomes negative on %s",
						acc.Name,
						acc.MinDate.Format("January 2, 2006"),
					),
					fmt.Errorf("case %d, balance check failed: %w", c, ledger.ErrNegativeBalance),
					"appledger.UpdateExpense",
				)
				return nil, err
			}
		}
	}

	if err := rtx.Update(ctx, t); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	// Monthly summary bookkeeping.
	newAccountID := string(fromAccID)

	if accountChanged {
		// Case A: account changed.
		// PREV account: remove PREV expense from PREV month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			0,
			-prevAmount,
			prevAccountID,
			monthlysummary.NewMonth(prevDate),
		); err != nil {
			return nil, err
		}
		// NEW account: add NEW expense at NEW month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			0,
			amount,
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
			// Different months: remove PREV expense from PREV month, add NEW expense to NEW month.
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				0,
				-prevAmount,
				newAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				0,
				amount,
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
				0,
				amount-prevAmount,
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
			"appledger.UpdateExpense",
		)
		return nil, err
	}

	return t, nil
}
