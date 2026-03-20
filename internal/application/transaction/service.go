package transaction

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

const LIMIT = 30

type Service struct {
	repo transaction.Repository
}

func NewService(repo transaction.Repository) *Service {
	return &Service{repo: repo}
}

// create income
// create expense
// create transfer

func (s *Service) ReadOneById(ctx context.Context, id string) (*transaction.Transaction, error) {
	newId, err := transaction.NewId(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindById(ctx, newId)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
	page uint,
) ([]*transaction.Transaction, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := transaction.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.repo.FindAllByUser(ctx, id, LIMIT, LIMIT*page-LIMIT)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) ReadAllByCategory(
	ctx context.Context,
	cid string,
	page uint,
) ([]*transaction.Transaction, error) {
	newId, err := transaction.NewId(cid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByCategory(ctx, newId, LIMIT, LIMIT*page-LIMIT)
}

func (s *Service) ReadWithFilters(
	ctx context.Context,
	filters *transaction.Filters,
	page uint,
) ([]*transaction.Transaction, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := transaction.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}
	return s.repo.FindWithFilters(ctx, id, filters, LIMIT, LIMIT*page-LIMIT)
}
