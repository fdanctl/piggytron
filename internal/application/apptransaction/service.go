package apptransaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

type Service struct {
	repo transaction.Repository
	db   *sql.DB
}

func NewService(r transaction.Repository, db *sql.DB) *Service {
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
) (*transaction.Transaction, error) {
	uid, err := util.ParseID[transaction.ID](userID)
	if err != nil {
		return nil, err
	}

	cid, err := util.ParseID[transaction.ID](categoryID)
	if err != nil {
		return nil, err
	}

	id, err := util.NewID[transaction.ID]()
	if err != nil {
		return nil, err
	}

	toAccID, err := util.ParseID[transaction.ID](dstAccID)
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewTransactionRepository(tx)

	acc, err := qtx.FindWithSum(ctx, dstAccID)
	if err != nil {
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
		return nil, err
	}

	t, err := transaction.NewIncome(
		id,
		uid,
		toAccID,
		cid,
		amount,
		description,
		date,
	)
	if err != nil {
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
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
) (*transaction.Transaction, error) {
	uid, err := util.ParseID[transaction.ID](userID)
	if err != nil {
		return nil, err
	}

	cid, err := util.ParseID[transaction.ID](catID)
	if err != nil {
		return nil, err
	}

	fromAccID, err := util.ParseID[transaction.ID](srcAccID)
	if err != nil {
		return nil, err
	}

	id, err := util.NewID[transaction.ID]()
	if err != nil {
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewTransactionRepository(tx)

	acc, err := qtx.FindWithSum(ctx, srcAccID)
	if err != nil {
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
		return nil, err
	}

	t, err := transaction.NewExpense(
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
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
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
) (*transaction.Transaction, error) {
	uid, err := util.ParseID[transaction.ID](userID)
	if err != nil {
		return nil, err
	}

	id, err := util.NewID[transaction.ID]()
	if err != nil {
		return nil, err
	}

	fromAccID, err := util.ParseID[transaction.ID](srcAccID)
	if err != nil {
		return nil, err
	}

	toAccID, err := util.ParseID[transaction.ID](dstAccID)
	if err != nil {
		return nil, err
	}

	var cid *transaction.ID
	if categoryID != "" {
		tempID, err := util.ParseID[transaction.ID](categoryID)
		if err != nil {
			return nil, err
		}
		cid = &tempID
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewTransactionRepository(tx)

	fromAccount, err := qtx.FindWithSum(ctx, srcAccID)
	if err != nil {
		return nil, err
	}
	toAccount, err := qtx.FindWithSum(ctx, dstAccID)
	if err != nil {
		return nil, err
	}

	var accCID *transaction.ID
	if toAccount.Category.ID != util.ZeroUUID {
		temp := transaction.ID(toAccount.Category.ID)
		accCID = &temp
	}

	t, err := transaction.NewTransfer(
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
		toAccount.Type,
		toAccount.IsSaving != nil && *toAccount.IsSaving,
	)
	if err != nil {
		return nil, err
	}

	err = rtx.Create(ctx, t)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

// Update

func (s *Service) Delete(ctx context.Context, id string) error {
	tid, err := util.ParseID[transaction.ID](id)
	if err != nil {
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewTransactionRepository(tx)

	t, err := rtx.FindByID(ctx, tid)
	if err != nil {
		return err
	}

	var toAccBalance *int
	if t.ToAccountID() != nil {
		toacc, err := qtx.FindWithSum(ctx, string(*t.ToAccountID()))
		if err != nil {
			return err
		}
		toAccBalance = &toacc.Sum
	}

	if err = t.CanBeDeleted(toAccBalance); err != nil {
		return err
	}

	if err = rtx.Delete(ctx, t.ID()); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Service) FindOneByID(ctx context.Context, id string) (*transaction.Transaction, error) {
	tid, err := util.ParseID[transaction.ID](id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, tid)
}

func (s *Service) FindAllByUser(
	ctx context.Context,
	userID string,
	page uint,
) ([]*transaction.Transaction, error) {
	uid, err := util.ParseID[transaction.ID](userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) FindAllByAccount(
	ctx context.Context,
	accountID string,
) ([]*transaction.Transaction, error) {
	aid, err := util.ParseID[transaction.ID](accountID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindAllByAccount(ctx, aid)
}

func (s *Service) FindAllByCategory(
	ctx context.Context,
	categoryID string,
) ([]*transaction.Transaction, error) {
	cid, err := util.ParseID[transaction.ID](categoryID)
	if err != nil {
		return nil, err
	}

	return s.repo.FindAllByCategory(ctx, cid)
}
