package account

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/google/uuid"
)

type Service struct {
	repo account.Repository
}

var (
	ErrDuplicate = errors.New("duplicate bank")
	ErrNotNumber = errors.New("not number")
)

func NewService(repo account.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBank(
	ctx context.Context,
	userID string,
	name string,
	currency string,
) (*account.Account, error) {
	uid, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindBankByNameAndUser(ctx, uid, name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := account.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	account, err := account.NewBank(id, uid, name, currency)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) CreateGoal(
	ctx context.Context,
	userID string,
	name string,
	currency string,
	targetAmount string,
	targetDate string,
	categoryID string,
) (*account.Account, error) {
	uid, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	_, err = s.repo.FindGoalByNameAndUser(ctx, account.ID(userID), name)
	if err == nil {
		return nil, ErrDuplicate
	}

	id, err := account.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}
	cid, err := account.NewID(categoryID)
	if err != nil {
		return nil, err
	}

	tAmount, err := strconv.Atoi(targetAmount)
	if err != nil {
		return nil, ErrNotNumber
	}

	tDate, err := time.Parse(time.DateOnly, targetDate)
	var pDate *time.Time
	if err == nil {
		pDate = &tDate
	}
	account, err := account.NewGoal(id, uid, name, currency, tAmount, pDate, cid)
	if err != nil {
		return nil, err
	}

	err = s.repo.Save(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) ReadOneByID(ctx context.Context, id string) (*account.Account, error) {
	newID, err := account.NewID(id)
	if err != nil {
		return nil, err
	}
	return s.repo.FindByID(ctx, newID)
}

func (s *Service) ReadAllByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	id, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) ReadAllBanksByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	id, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllBanksByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *Service) ReadAllGoalsByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	id, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	accounts, err := s.repo.FindAllGoalsByUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return accounts, nil
}
