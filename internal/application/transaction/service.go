package transaction

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
)

type Service struct {
	repo transaction.Repository
}

func NewService(r transaction.Repository) *Service {
	return &Service{repo: r}
}

// create income
// create expense
// create transfer

func (s *Service) ReadOneByID(ctx context.Context, id string) (*transaction.Transaction, error) {
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
	newID, err := transaction.NewID(cid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByCategory(ctx, newID)
}
