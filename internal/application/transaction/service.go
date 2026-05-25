package transaction

import (
	"context"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/google/uuid"
)

type Service struct {
	repo transaction.Repository
}

func NewService(r transaction.Repository) *Service {
	return &Service{repo: r}
}

// create expense
// create transfer

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

	transaction, err := transaction.NewIncome(
		id, uid, toAccID, cid, amount, description, date)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
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

	transaction, err := transaction.NewExpense(
		id, uid, fromAccID, cid, amount, description, date)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
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

	transaction, err := transaction.NewTransfer(
		id, uid, fromAccID, toAccID, cid, amount, description, date)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}
	return transaction, nil
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
