package bank

import (
	"context"
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/bank"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/google/uuid"
)

type Service struct {
	repo bank.Repository
}

var ErrDuplicate = errors.New("duplicate bank")

func NewService(repo bank.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBank(
	ctx context.Context,
	name string,
	currency string,
) (*bank.Bank, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, errors.New("nil context")
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, errors.New("not sessionInfo")
	}

	_, err := s.repo.FindByNameAndUser(ctx, bank.ID(sessionInfo.UserId), name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := bank.NewId(uuid.New().String())
	if err != nil {
		return nil, err
	}
	uid, err := bank.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	bank, err := bank.New(id, uid, name, currency)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, bank)
	if err != nil {
		return nil, err
	}
	return bank, nil
}

func (s *Service) ReadOneById(ctx context.Context, id string) (*bank.Bank, error) {
	newId, err := bank.NewId(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindById(ctx, newId)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
) ([]*bank.Bank, error) {
	v := ctx.Value("user")
	if v == nil {
		return nil, nil
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, nil
	}

	id, err := bank.NewId(sessionInfo.UserId)
	if err != nil {
		return nil, err
	}

	banks, err := s.repo.FindAllByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return banks, nil
}
