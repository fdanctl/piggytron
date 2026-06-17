package appbudget

import (
	"context"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/util"
)

type Service struct {
	repo budget.Repository
}

var (
	ErrDuplicate      = errors.New("duplicate budget")
	ErrNegativeNumber = errors.New("negative number")
)

func NewService(r budget.Repository) *Service {
	return &Service{repo: r}
}

func (s *Service) CreateBudget(
	ctx context.Context,
	userID string,
	categoryID string,
	month time.Time,
	amount int,
) (*budget.Budget, error) {
	uid, err := util.ParseID[budget.ID](userID)
	if err != nil {
		return nil, err
	}

	cid, err := util.ParseID[budget.ID](categoryID)
	if err != nil {
		return nil, err
	}

	// TODO handle duplicate error

	id, err := util.NewID[budget.ID]()
	if err != nil {
		return nil, err
	}

	budget, err := budget.New(
		id,
		uid,
		cid,
		month,
		amount,
	)
	if err != nil {
		return nil, err
	}

	err = s.repo.Create(ctx, budget)
	if err != nil {
		return nil, err
	}
	return budget, nil
}

func (s *Service) UpdateBudgetAmount(
	ctx context.Context,
	id string,
	amount int,
) error {
	bid, err := util.ParseID[budget.ID](id)
	if err != nil {
		return err
	}

	if amount < 0 {
		return ErrNegativeNumber
	}

	return s.repo.UpdateAmount(ctx, bid, amount)
}

func (s *Service) FindBudget(
	ctx context.Context,
	id string,
) (*budget.Budget, error) {
	bid, err := util.ParseID[budget.ID](id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, bid)
}
