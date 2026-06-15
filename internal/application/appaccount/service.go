package appaccount

import (
	"context"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/google/uuid"
)

type Service struct {
	repo account.Repository
}

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidDate   = errors.New("invalid date")
)

func NewService(repo account.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindOneByID(ctx context.Context, id string) (*account.Account, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	accID, err := account.NewID(id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, accID)
}

func (s *Service) CreateBank(
	ctx context.Context,
	userID string,
	name string,
	currency string,
	isSaving bool,
) (*account.Account, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	uid, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := account.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	account, err := account.NewBank(id, uid, name, currency, isSaving)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, account)
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
	targetAmount int,
	startDate time.Time,
	targetDate *time.Time,
	categoryID string,
) (*account.Account, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	uid, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	id, err := account.NewID(uuid.New().String())
	if err != nil {
		return nil, err
	}
	cid, err := account.NewID(categoryID)
	if err != nil {
		return nil, err
	}

	account, err := account.NewGoal(
		id,
		uid,
		name,
		currency,
		targetAmount,
		startDate,
		targetDate,
		cid,
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (s *Service) UpdateGoal(
	ctx context.Context,
	id string,
	userID string,
	name string,
	currency string,
	targetAmount int,
	startDate time.Time,
	targetDate *time.Time,
	categoryID string,
) (*account.Account, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}

	uid, err := account.NewID(userID)
	if err != nil {
		return nil, err
	}

	cid, err := account.NewID(categoryID)
	if err != nil {
		return nil, err
	}

	goal, err := s.repo.FindByID(ctx, account.ID(id))
	if err != nil {
		return nil, err
	}
	if goal.UserID() != uid {
		return nil, errors.New("not found")
	}

	// update
	goal.ChangeName(name)
	goal.ChangeTargetAmount(targetAmount)
	goal.ChangeStartDate(startDate)
	goal.ChangeTargetDate(targetDate)
	goal.ChangeCategory(cid)

	err = s.repo.Update(ctx, goal)
	if err != nil {
		return nil, err
	}

	return goal, nil
}

func (s *Service) FindAllByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

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

func (s *Service) FindAllBanksByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

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

func (s *Service) FindAllGoalsByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

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
