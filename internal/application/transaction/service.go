package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/google/uuid"
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
	catID string,
	dstID string,
) (*transaction.Transaction, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(catID)
	if err != nil {
		return nil, err
	}

	uid, err := transaction.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := transaction.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	toAccID, err := transaction.NewID(dstID)
	if err != nil {
		return nil, err
	}

	cid, err := transaction.NewID(catID)
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

	acc, err := qtx.FindWithSum(ctx, dstID)
	if err != nil {
		return nil, err
	}

	var accCID *account.ID
	if acc.Category.ID != "00000000-0000-0000-0000-000000000000" {
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
	srcID string,
) (*transaction.Transaction, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(catID)
	if err != nil {
		return nil, err
	}

	uid, err := transaction.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := transaction.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	fromAccID, err := transaction.NewID(srcID)
	if err != nil {
		return nil, err
	}

	cid, err := transaction.NewID(catID)
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

	acc, err := qtx.FindWithSum(ctx, srcID)
	if err != nil {
		return nil, err
	}

	var accCID *account.ID
	if acc.Category.ID != "00000000-0000-0000-0000-000000000000" {
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
	catID string,
	srcID string,
	dstID string,
) (*transaction.Transaction, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	uid, err := transaction.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := transaction.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	fromAccID, err := transaction.NewID(srcID)
	if err != nil {
		return nil, err
	}

	toAccID, err := transaction.NewID(dstID)
	if err != nil {
		return nil, err
	}

	var cid *transaction.ID
	if catID != "" {
		tempID, err := transaction.NewID(catID)
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

	fromAccount, err := qtx.FindWithSum(ctx, srcID)
	if err != nil {
		return nil, err
	}
	toAccount, err := qtx.FindWithSum(ctx, dstID)
	if err != nil {
		return nil, err
	}

	var accCID *transaction.ID
	if toAccount.Category.ID != "00000000-0000-0000-0000-000000000000" {
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
	_, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	ID, err := transaction.NewID(id)
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

	t, err := rtx.FindByID(ctx, ID)
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

func (s *Service) ReadOneByID(ctx context.Context, id string) (*transaction.Transaction, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	newID, err := transaction.NewID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, newID)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
	userID string,
	page uint,
) ([]*transaction.Transaction, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	uid, err := transaction.NewID(userID)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) ReadAllByAccount(
	ctx context.Context,
	aid string,
) ([]*transaction.Transaction, error) {
	_, err := uuid.Parse(aid)
	if err != nil {
		return nil, err
	}

	newID, err := transaction.NewID(aid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByAccount(ctx, newID)
}

func (s *Service) ReadAllByCategory(
	ctx context.Context,
	cid string,
) ([]*transaction.Transaction, error) {
	_, err := uuid.Parse(cid)
	if err != nil {
		return nil, err
	}

	newID, err := transaction.NewID(cid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByCategory(ctx, newID)
}
