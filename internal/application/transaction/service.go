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

	transactions, err := s.repo.FindAllByUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return transactions, nil
}

func (s *Service) ReadAllByCategory(
	ctx context.Context,
	cid string,
) ([]*transaction.Transaction, error) {
	newId, err := transaction.NewId(cid)
	if err != nil {
		return nil, err
	}
	return s.repo.FindAllByCategory(ctx, newId)
}

func (s *Service) ReadWithFilters(
	ctx context.Context,
	filters *transaction.Filters,
	page uint,
) ([]*transaction.Transaction, bool, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, false, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, false, nil
	}

	id, err := transaction.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, false, err
	}

	transactions, err := s.repo.FindWithFilters(ctx, id, filters, LIMIT+1, LIMIT*page-LIMIT)
	if err != nil {
		return nil, false, err
	}

	var hasMore bool
	if len(transactions) == LIMIT+1 {
		hasMore = true
		transactions = transactions[0 : len(transactions)-1]
	}

	return transactions, hasMore, nil
}

func (s *Service) CountFilteredResults(
	ctx context.Context,
	filters *transaction.Filters,
) (uint, error) {
	v := ctx.Value("user")
	if v == nil {
		return 0, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return 0, nil
	}

	id, err := transaction.NewId(sessionInfo.UserId)
	if err != nil {
		return 0, err
	}
	return s.repo.CountFilteredResults(ctx, id, filters)
}
