package appledger

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/internal/util"
)

type Service struct {
	repo ledger.Repository
	db   *sql.DB
}

func NewService(r ledger.Repository, db *sql.DB) *Service {
	return &Service{repo: r, db: db}
}

func (s *Service) CreateIncome(
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
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	cid, err := util.ParseID[ledger.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", cid),
			fmt.Errorf("failed parsing id '%s': %w", cid, err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	id, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	toAccID, err := util.ParseID[ledger.ID](dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", dstAccID),
			fmt.Errorf("failed parsing id '%s': %w", dstAccID, err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating transaction: %w", err),
			"appledger.CreateIncome",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)

	acc, err := qtx.FindWithSum(ctx, dstAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Source account not found",
			fmt.Errorf("failed to find account '%s': %w", dstAccID, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	var accCID *account.ID
	if acc.Category.ID != util.ZeroUUID {
		temp := account.ID(acc.ID)
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

	if err := a.CanReceiveIncome(); err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s of type %s can't receive income ledger entries", a.Name(), a.Type()),
			fmt.Errorf("%s of type %s can't receive income: %w", a.Name(), a.Type(), err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	t, err := ledger.NewIncome(
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
			"Failed to create income",
			fmt.Errorf("failed to create income: %w", err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appledger.CreateIncome",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CreateIncome",
		)
		return nil, err
	}

	return t, nil
}

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

	acc, err := qtx.FindWithSum(ctx, srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Source account not found",
			fmt.Errorf("failed to find account '%s': %w", srcAccID, err),
			"appledger.CreateExpense",
		)
		return nil, err
	}

	var accCID *account.ID
	if acc.Category.ID != util.ZeroUUID {
		temp := account.ID(acc.ID)
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

	if err := a.CanReceiveIncome(); err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			fmt.Sprintf("%s of type %s can't receive income ledger entries", a.Name(), a.Type()),
			fmt.Errorf("%s of type %s can't receive income: %w", a.Name(), a.Type(), err),
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
		acc.Sum,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create expense",
			fmt.Errorf("failed to create expense: %w", err),
			"appledger.CreateExpense",
		)
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

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CreateExpense",
		)
		return nil, err
	}

	return t, nil
}

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
				fmt.Sprintf("%s is not a valid id", *cid),
				fmt.Errorf("failed parsing id '%s': %w", *cid, err),
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

	fromAccount, err := qtx.FindWithSum(ctx, srcAccID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Source account not found",
			fmt.Errorf("failed to find account '%s': %w", srcAccID, err),
			"appledger.CreateTransfer",
		)
		return nil, err
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
	if toAccount.IsSaving != nil && *toAccount.IsSaving {
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
		fromAccount.Sum,
		accCID,
		toAccountCatType,
		toAccount.IsSaving != nil && *toAccount.IsSaving,
	)
	if err != nil {
		msg := "Failed to create transfer"
		if errors.Is(err, ledger.ErrNegativeBalance) {
			msg = fmt.Sprintf("%s account becomes negative", fromAccount.Name)
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

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appledger.CreateTransfer",
		)
		return nil, err
	}

	return t, nil
}

// Update

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

	t, err := rtx.FindByID(ctx, tid)
	if err != nil {
		return err
	}

	var toAccWithSum *query.AccountWithSum
	if t.ToAccountID() != nil {
		toacc, err := qtx.FindWithSum(ctx, string(*t.ToAccountID()))
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Destination account not found",
				fmt.Errorf("failed to find account '%s': %w", *t.ToAccountID(), err),
				"appledger.Delete",
			)
			return err
		}
		toAccWithSum = toacc
	}

	if err = t.CanBeDeleted(&toAccWithSum.Sum); err != nil {
		if errors.Is(err, ledger.ErrNegativeBalance) {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s account becomes negative", toAccWithSum.Name),
				fmt.Errorf("%s account becomes negative: %w", *t.ToAccountID(), err),
				"appledger.Delete",
			)
		}
		return err
	}

	if err = rtx.Delete(ctx, t.ID()); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to delete transaction': %w", err),
			"appledger.Delete",
		)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appledger.Delete",
		)
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
