package transaction

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	t, err := transaction.NewIncome(
		id, uid, toAccID, cid, amount, description, date)
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

	account, err := qtx.FindWithSum(ctx, dstID)
	if err != nil {
		return nil, err
	}

	if account.IsSaving != nil && *account.IsSaving {
		return nil, transaction.ErrInvalidAccount
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

	t, err := transaction.NewExpense(
		id, uid, fromAccID, cid, amount, description, date)
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

	account, err := qtx.FindWithSum(ctx, srcID)
	if err != nil {
		return nil, err
	}

	if account.Sum-amount < 0 {
		return nil, transaction.ErrNegativeBalance
	}

	if account.IsSaving != nil && !*account.IsSaving {
		return nil, transaction.ErrInvalidAccount
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

	t, err := transaction.NewTransfer(
		id, uid, fromAccID, toAccID, cid, amount, description, date)
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
	cqtx := postgres.NewCategoryQueryService(tx)

	fromAccount, err := qtx.FindWithSum(ctx, srcID)
	toAccount, err := qtx.FindWithSum(ctx, dstID)
	if err != nil {
		return nil, err
	}

	if fromAccount.Sum-amount < 0 {
		return nil, transaction.ErrNegativeBalance
	}

	if toAccount.Type == "goal" && toAccount.Category.ID != catID {
		return nil, fmt.Errorf(
			"%w: %s goal must be %s category",
			transaction.ErrGoalCategory,
			toAccount.Name,
			toAccount.Category.Name,
		)
	}

	if toAccount.IsSaving != nil && *toAccount.IsSaving {
		if catID == "" {
			return nil, transaction.ErrNotSavingsCategory
		}
		cat, err := cqtx.FindByID(ctx, catID)
		if err != nil {
			return nil, err
		}
		if cat.Type != "savings" {
			return nil, transaction.ErrNotSavingsCategory
		}
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
