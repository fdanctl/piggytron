package transaction

import (
	"context"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

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

func (s *Service) RealAllByUser(
	ctx context.Context,
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
