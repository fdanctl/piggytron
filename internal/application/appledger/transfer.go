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

// CreateTransfer moves money between two accounts, enforcing the goal and
// savings category rules and the source's solvency, and updates both monthly
// summaries.
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
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
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
	var fromAccountCID *account.ID
	if fromAccount.Category.ID != util.ZeroUUID {
		temp := account.ID(fromAccount.Category.ID)
		fromAccountCID = &temp
	}
	fromA := account.Rehydrate(
		account.ID(fromAccount.ID),
		account.ID(fromAccount.UserID),
		account.AccountType(fromAccount.Type),
		fromAccount.Name,
		account.AccountStatus(fromAccount.Status),
		fromAccount.TargetAmount,
		fromAccount.StartDate,
		fromAccount.TargetDate,
		fromAccountCID,
		fromAccount.Currency,
		fromAccount.ClosedAt,
		fromAccount.CreatedAt,
		fromAccount.UpdatedAt,
	)

	if fromA.IsClosed() {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Can't change closed accounts",
			account.ErrClosedAccount,
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
	var toAccountCID *account.ID
	if toAccount.Category.ID != util.ZeroUUID {
		temp := account.ID(toAccount.Category.ID)
		toAccountCID = &temp
	}
	toA := account.Rehydrate(
		account.ID(toAccount.ID),
		account.ID(toAccount.UserID),
		account.AccountType(toAccount.Type),
		toAccount.Name,
		account.AccountStatus(toAccount.Status),
		toAccount.TargetAmount,
		toAccount.StartDate,
		toAccount.TargetDate,
		toAccountCID,
		toAccount.Currency,
		toAccount.ClosedAt,
		toAccount.CreatedAt,
		toAccount.UpdatedAt,
	)

	if toA.IsClosed() {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Can't change closed accounts",
			account.ErrClosedAccount,
			"appledger.CreateTransfer",
		)
	}

	var accCID *ledger.ID
	if toAccount.Category.ID != util.ZeroUUID {
		temp := ledger.ID(toAccount.Category.ID)
		accCID = &temp
	}

	var toAccountCatType string
	if toAccount.Type == string(account.SavingsType) && categoryID != "" {
		cattx := postgres.NewCategoryQueryService(tx)
		cat, err := cattx.FindByID(ctx, categoryID)
		if err != nil {
			return nil, errs.NewInternalAppError(err, "appledger.CreateTransfer")
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
		toAccount.Type == string(account.SavingsType),
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
				"Transfers to %s must be %s category",
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
			fmt.Errorf("failed to create transfer: %w", err),
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

// UpdateTransfer modifies a transfer entry, re-running the goal/savings and
// solvency checks affected by the change and reconciling monthly summaries on
// both sides.
func (s *Service) UpdateTransfer(
	ctx context.Context,
	id string,
	userID string,
	amount int,
	currency string,
	description string,
	date time.Time,
	categoryID string,
	srcAccID string,
	dstAccID string,
) (*ledger.Entry, error) {
	tid, err := util.ParseID[ledger.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appledger.UpdateTransfer",
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
				"appledger.UpdateTransfer",
			)
			return nil, err
		}
		cid = &tempID
	}

	fromAccID, err := util.ParseID[ledger.ID](srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", srcAccID),
			fmt.Errorf("failed parsing id '%s': %w", srcAccID, err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	toAccID, err := util.ParseID[ledger.ID](dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", dstAccID),
			fmt.Errorf("failed parsing id '%s': %w", dstAccID, err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)

	t, err := rtx.FindByID(ctx, tid)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find transaction",
				fmt.Errorf("failed to find transaction '%s': %w", tid, err),
				"appledger.UpdateTransfer",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find transaction '%s': %w", id, err),
				"appledger.UpdateTransfer",
			)
		}
		return nil, err
	}

	if t.UserID() != uid || t.Type() != "transfer" {
		return nil, ledger.ErrNotFound
	}

	fromAccChanged := t.FromAccountID() == nil || string(*t.FromAccountID()) != string(fromAccID)
	toAccChanged := t.ToAccountID() == nil || string(*t.ToAccountID()) != string(toAccID)
	prevAmount := t.Amount()
	prevDate := t.Date()
	prevFromAccountID := string(*t.FromAccountID())
	prevToAccountID := string(*t.ToAccountID())

	// Update
	if err := t.UpdateTransfer(fromAccID, toAccID, amount, description, date); err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to update transfer",
			fmt.Errorf("failed to update transfer: %w", err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	// verify if the account is closed
	oldToAcc, err := qtx.FindWithSum(ctx, prevToAccountID)
	if err != nil {
		return nil, errs.NewInternalAppError(err, "appledger.UpdateTransfer")
	}
	var oldToAccCID *account.ID
	if oldToAcc.Category.ID != util.ZeroUUID {
		temp := account.ID(oldToAcc.Category.ID)
		oldToAccCID = &temp
	}
	oldToA := account.Rehydrate(
		account.ID(oldToAcc.ID),
		account.ID(oldToAcc.UserID),
		account.AccountType(oldToAcc.Type),
		oldToAcc.Name,
		account.AccountStatus(oldToAcc.Status),
		oldToAcc.TargetAmount,
		oldToAcc.StartDate,
		oldToAcc.TargetDate,
		oldToAccCID,
		oldToAcc.Currency,
		oldToAcc.ClosedAt,
		oldToAcc.CreatedAt,
		oldToAcc.UpdatedAt,
	)

	if oldToA.IsClosed() {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Can't change closed accounts",
			account.ErrClosedAccount,
			"appledger.UpdateTransfer",
		)
	}

	oldFromAcc, err := qtx.FindWithSum(ctx, prevFromAccountID)
	if err != nil {
		return nil, errs.NewInternalAppError(err, "appledger.UpdateTransfer")
	}
	var oldFromAccCID *account.ID
	if oldFromAcc.Category.ID != util.ZeroUUID {
		temp := account.ID(oldFromAcc.Category.ID)
		oldFromAccCID = &temp
	}
	oldFromA := account.Rehydrate(
		account.ID(oldFromAcc.ID),
		account.ID(oldFromAcc.UserID),
		account.AccountType(oldFromAcc.Type),
		oldFromAcc.Name,
		account.AccountStatus(oldFromAcc.Status),
		oldFromAcc.TargetAmount,
		oldFromAcc.StartDate,
		oldFromAcc.TargetDate,
		oldFromAccCID,
		oldFromAcc.Currency,
		oldFromAcc.ClosedAt,
		oldFromAcc.CreatedAt,
		oldFromAcc.UpdatedAt,
	)

	if oldFromA.IsClosed() {
		return nil, errs.NewAppError(
			errs.KindBusinessRule,
			"Can't change closed accounts",
			account.ErrClosedAccount,
			"appledger.UpdateTransfer",
		)
	}

	if toAccChanged || cid != t.ExpenseCategoryID() {
		toAccount, err := qtx.FindWithSum(ctx, dstAccID)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Destination account not found",
				fmt.Errorf("failed to find account '%s': %w", dstAccID, err),
				"appledger.UpdateTransfer",
			)
			return nil, err
		}

		var accCID *ledger.ID
		if toAccount.Category.ID != util.ZeroUUID {
			temp := ledger.ID(toAccount.Category.ID)
			accCID = &temp
		}

		a := account.Rehydrate(
			account.ID(toAccount.ID),
			account.ID(toAccount.UserID),
			account.AccountType(toAccount.Type),
			toAccount.Name,
			account.AccountStatus(toAccount.Status),
			toAccount.TargetAmount,
			toAccount.StartDate,
			toAccount.TargetDate,
			oldToAccCID,
			toAccount.Currency,
			toAccount.ClosedAt,
			toAccount.CreatedAt,
			toAccount.UpdatedAt,
		)

		if a.IsClosed() {
			return nil, errs.NewAppError(
				errs.KindBusinessRule,
				"Can't change closed accounts",
				account.ErrClosedAccount,
				"appledger.UpdateTransfer",
			)
		}

		var toAccountCatType string
		if toAccount.Type == string(account.SavingsType) && categoryID != "" {
			cattx := postgres.NewCategoryQueryService(tx)
			cat, err := cattx.FindByID(ctx, categoryID)
			if err != nil {
				return nil, errs.NewInternalAppError(err, "appledger.UpdateTransfer")
			}
			toAccountCatType = cat.Type
		}

		if err := t.UpdateTransferToAccountAndCategory(
			cid,
			toAccID,
			accCID,
			toAccountCatType,
			toAccount.Type == string(account.SavingsType),
		); err != nil {
			msg := "Failed to update transfer"
			if errors.Is(err, ledger.ErrGoalCategory) {
				msg = fmt.Sprintf(
					"Transfers to %s must be %s category",
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
				fmt.Errorf("failed to update transfer: %w", err),
				"appledger.UpdateTransfer",
			)
			return nil, err
		}
	}

	// Balance verification
	if toAccChanged {
		// Verify the PREV account (it loses this income)
		acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevToAccountID, prevDate, nil, &id)
		if err != nil {
			return nil, errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", prevToAccountID, err),
				"appledger.UpdateTransfer",
			)
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
			oldToAccCID,
			acc.Currency,
			acc.ClosedAt,
			acc.CreatedAt,
			acc.UpdatedAt,
		)

		if a.IsClosed() {
			return nil, errs.NewAppError(
				errs.KindBusinessRule,
				"Can't change closed accounts",
				account.ErrClosedAccount,
				"appledger.UpdateExpense",
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
				"appledger.UpdateTransfer",
			)
			return nil, err
		}
	} else {
		// Same account
		c, err := getCaseFromTable(date.Compare(prevDate), cmp.Compare(amount, prevAmount))
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("don't know how could I make this error: %w", err),
				"appledger.UpdateTransfer",
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
			acc, err := qtx.GetAccountWithMinRunningBalance(
				ctx,
				prevToAccountID,
				prevDate,
				nil,
				&id,
			)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevToAccountID, err),
					"appledger.UpdateTransfer",
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
					"appledger.UpdateTransfer",
				)
				return nil, err
			}

		case 7, 8, 9:
			// Check from PREV date to NEW date (not included) the balance without the PREV amount
			until := time.Date(date.Year(), date.Month(), date.Day()-1, 0, 0, 0, 0, date.Location())
			acc, err := qtx.GetAccountWithMinRunningBalance(
				ctx,
				prevToAccountID,
				prevDate,
				&until,
				&id,
			)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevToAccountID, err),
					"appledger.UpdateTransfer",
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
					"appledger.UpdateTransfer",
				)
				return nil, err
			}
			if c == 7 {
				// Check from NEW date the balance without the PREV amount + the NEW amount
				acc, err := qtx.GetAccountWithMinRunningBalance(
					ctx,
					prevToAccountID,
					date,
					nil,
					&id,
				)
				if err != nil {
					return nil, errs.NewInternalAppError(
						fmt.Errorf("failed to find account '%s': %w", prevToAccountID, err),
						"appledger.UpdateTransfer",
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
						"appledger.UpdateTransfer",
					)
					return nil, err
				}
			}
		}
	}

	// Account changed: NEW account gains expense → verify NEW account stays solvent.
	if fromAccChanged {
		// Verify the NEW account (it gains expense)
		acc, err := qtx.GetAccountWithMinRunningBalance(ctx, srcAccID, date, nil, nil)
		if err != nil {
			return nil, errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", prevFromAccountID, err),
				"appledger.UpdateTransfer",
			)
		}

		delta := -amount
		if srcAccID == prevToAccountID {
			// the NEW now withdraws money from an account the recieved money in PREV
			delta -= prevAmount
		}

		if acc.MinRunningBalance+delta < 0 {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf(
					"%s becomes negative on %s",
					acc.Name,
					acc.MinDate.Format("January 2, 2006"),
				),
				fmt.Errorf("balance check failed: %w", ledger.ErrNegativeBalance),
				"appledger.UpdateTransfer",
			)
			return nil, err
		}
	} else {
		// Same account
		c, err := getCaseFromTable(date.Compare(prevDate), cmp.Compare(amount, prevAmount))
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("don't know how could I make this error: %w", err),
				"appledger.UpdateTransfer",
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
					fmt.Errorf("failed to find account '%s': %w", prevFromAccountID, err),
					"appledger.UpdateTransfer",
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
					"appledger.UpdateTransfer",
				)
				return nil, err
			}
			if c == 3 {
				// Check from PREV date the balance + PREV amount - NEW amount >= 0
				acc, err := qtx.GetAccountWithMinRunningBalance(
					ctx,
					prevFromAccountID,
					prevDate,
					nil,
					nil,
				)
				if err != nil {
					return nil, errs.NewInternalAppError(
						fmt.Errorf("failed to find account '%s': %w", prevFromAccountID, err),
						"appledger.UpdateTransfer",
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
						"appledger.UpdateTransfer",
					)
					return nil, err
				}
			}

		case 6, 9:
			// Check from NEW date the balance + PREV amount - NEW amount >= 0
			acc, err := qtx.GetAccountWithMinRunningBalance(ctx, prevFromAccountID, date, nil, nil)
			if err != nil {
				return nil, errs.NewInternalAppError(
					fmt.Errorf("failed to find account '%s': %w", prevFromAccountID, err),
					"appledger.UpdateTransfer",
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
					"appledger.UpdateTransfer",
				)
				return nil, err
			}
		}
	}

	if err := rtx.Update(ctx, t); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	// Monthly summary bookkeeping.
	newToAccountID := string(toAccID)

	if toAccChanged {
		// Case A: account changed.
		// PREV account: remove PREV income from PREV month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			-prevAmount,
			0,
			prevToAccountID,
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
			newToAccountID,
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
				newToAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				amount,
				0,
				newToAccountID,
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
				newToAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
		}
	}
	newFromAccountID := string(fromAccID)

	if fromAccChanged {
		// Case A: account changed.
		// PREV account: remove PREV expense from PREV month.
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			0,
			-prevAmount,
			prevFromAccountID,
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
			newFromAccountID,
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
				newFromAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
			if err := s.updateMonthlySummary(
				ctx,
				mstx,
				0,
				amount,
				newFromAccountID,
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
				newFromAccountID,
				prevMonth,
			); err != nil {
				return nil, err
			}
		}
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit: %w", err),
			"appledger.UpdateTransfer",
		)
		return nil, err
	}

	return t, nil
}
