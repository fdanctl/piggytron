package appledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

type Service struct {
	repo ledger.Repository
	db   *sql.DB
}

var ErrCanChangeType = errors.New("can't change transaction type")

func NewService(r ledger.Repository, db *sql.DB) *Service {
	return &Service{repo: r, db: db}
}

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
	}
	return nil, errors.New("invalid transaction type")
}

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
