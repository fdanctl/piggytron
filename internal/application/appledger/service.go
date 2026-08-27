// Package appledger implements the ledger use cases: creating, updating and
// deleting income, expense and transfer entries, and querying them. Every
// write keeps the affected accounts' monthly summaries in sync.
package appledger

import (
	"context"
	"database/sql"
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

// Service implements the ledger use cases. Every write runs inside a single
// database transaction and keeps the affected accounts' monthly summaries in
// sync.
type Service struct {
	repo ledger.Repository
	db   *sql.DB
}

// NewService wires the ledger service to its repository and database.
func NewService(r ledger.Repository, db *sql.DB) *Service {
	return &Service{repo: r, db: db}
}

// Update dispatches the request to the typed update matching ttype (income,
// expense or transfer).
func (s *Service) Update(
	ctx context.Context,
	ttype string,
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
	switch ttype {
	case "income":
		return s.UpdateIncome(
			ctx,
			id,
			userID,
			amount,
			currency,
			description,
			date,
			categoryID,
			dstAccID,
		)

	case "expense":
		return s.UpdateExpense(
			ctx,
			id,
			userID,
			amount,
			currency,
			description,
			date,
			categoryID,
			srcAccID,
		)

	case "transfer":
		return s.UpdateTransfer(
			ctx,
			id,
			userID,
			amount,
			currency,
			description,
			date,
			categoryID,
			srcAccID,
			dstAccID,
		)
	}
	return nil, errs.NewGenericBadRequestAppError(ledger.ErrInvalidType, "appledger.Update")
}

// Delete removes an entry, verifying the destination account stays solvent
// and subtracting the movement from both accounts' monthly summaries.
func (s *Service) Delete(ctx context.Context, id string) error {
	tid, err := util.ParseID[ledger.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.Delete",
		)
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.Delete",
		)
		return err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)

	t, err := rtx.FindByID(ctx, tid)
	if err != nil {
		return err
	}

	// if it's income or a transfer
	if t.ToAccountID() != nil {
		toAcc, err := qtx.GetAccountWithMinRunningBalance(
			ctx,
			string(*t.ToAccountID()),
			t.Date(),
			nil,
			nil,
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Destination account not found",
				fmt.Errorf("failed to find account '%s': %w", *t.ToAccountID(), err),
				"appledger.Delete",
			)
			return err
		}

		// verify if the account is closed
		oldAcc, err := qtx.FindWithSum(ctx, toAcc.ID)
		if err != nil {
			return errs.NewInternalAppError(err, "appledger.Delete")
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
			return errs.NewAppError(
				errs.KindBusinessRule,
				"Destination account is closed. Can't change closed accounts",
				account.ErrClosedAccount,
				"appledger.Delete",
			)
		}

		if err = t.CanBeDeleted(&toAcc.MinRunningBalance); err != nil {
			if errors.Is(err, ledger.ErrNegativeBalance) {
				err = errs.NewAppError(
					errs.KindValidation,
					fmt.Sprintf(
						"%s becomes negative on %s",
						toAcc.Name,
						toAcc.MinDate.Format("January 2, 2006"),
					),
					fmt.Errorf("%s account becomes negative: %w", *t.ToAccountID(), err),
					"appledger.Delete",
				)
			}
			return err
		}
	}

	// check from account is closed
	if t.FromAccountID() != nil {
		fromAcc, err := qtx.FindWithSum(ctx, string(*t.FromAccountID()))
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Source account not found",
				fmt.Errorf("failed to find account '%s': %w", *t.FromAccountID(), err),
				"appledger.Delete",
			)
			return err
		}

		// verify if the account is closed
		oldAcc, err := qtx.FindWithSum(ctx, fromAcc.ID)
		if err != nil {
			return errs.NewInternalAppError(err, "appledger.Delete")
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
			return errs.NewAppError(
				errs.KindBusinessRule,
				"Source account is closed. Can't change closed accounts",
				account.ErrClosedAccount,
				"appledger.Delete",
			)
		}
	}

	if err = rtx.Delete(ctx, t.ID()); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to delete transaction': %w", err),
			"appledger.Delete",
		)
		return err
	}

	// monthly summary — read current, subtract, write back
	month := monthlysummary.NewMonth(t.Date())

	if t.FromAccountID() != nil {
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			0,
			-t.Amount(),
			string(*t.FromAccountID()),
			month,
		); err != nil {
			return err
		}
	}

	if t.ToAccountID() != nil {
		if err := s.updateMonthlySummary(
			ctx,
			mstx,
			-t.Amount(),
			0,
			string(*t.ToAccountID()),
			month,
		); err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit: %w", err),
			"appledger.Delete",
		)
		return err
	}
	return nil
}

// FindOneByID returns a ledger entry by id.
func (s *Service) FindOneByID(ctx context.Context, id string) (*ledger.Entry, error) {
	tid, err := util.ParseID[ledger.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.FindOneByID",
		)
		return nil, err
	}

	t, err := s.repo.FindByID(ctx, tid)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"Failed to find transaction",
				fmt.Errorf("failed to find transaction '%s': %w", tid, err),
				"appledger.FindOneByID",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find transaction '%s': %w", id, err),
				"appledger.FindOneByID",
			)
		}
		return nil, err
	}

	return t, nil
}

// FindAllByUser returns the user's ledger entries.
func (s *Service) FindAllByUser(
	ctx context.Context,
	userID string,
	page uint,
) ([]*ledger.Entry, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appledger.FindAllByUser",
		)
		return nil, err
	}

	transactions, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to find user '%s' ledger entries: %w", uid, err),
			"appledger.FindAllByUser",
		)
		return nil, err
	}
	return transactions, nil
}

// FindAllByAccount returns the ledger entries of an account.
func (s *Service) FindAllByAccount(
	ctx context.Context,
	accountID string,
) ([]*ledger.Entry, error) {
	aid, err := util.ParseID[ledger.ID](accountID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", aid),
			fmt.Errorf("failed parsing id '%s': %w", aid, err),
			"appledger.FindAllByAccount",
		)
		return nil, err
	}

	transactions, err := s.repo.FindAllByAccount(ctx, aid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to find account's '%s' ledger entries: %w", aid, err),
			"appledger.FindAllByAccount",
		)
		return nil, err
	}
	return transactions, nil
}

// FindAllByCategory returns the ledger entries of a category.
func (s *Service) FindAllByCategory(
	ctx context.Context,
	categoryID string,
) ([]*ledger.Entry, error) {
	cid, err := util.ParseID[ledger.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", cid),
			fmt.Errorf("failed parsing id '%s': %w", cid, err),
			"appledger.FindAllByCategory",
		)
		return nil, err
	}

	transactions, err := s.repo.FindAllByCategory(ctx, cid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to find category '%s' ledger entries: %w", cid, err),
			"appledger.FindAllByCategory",
		)
		return nil, err
	}
	return transactions, nil
}
