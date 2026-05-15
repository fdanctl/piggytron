package budget

import (
	"context"
	"errors"
	"time"

	budget "github.com/fdanctl/piggytron/internal/domain/monthly_budget"
	"github.com/google/uuid"
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
	_, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	uid, err := budget.NewID(userID)
	if err != nil {
		return nil, err
	}

	_, err = uuid.Parse(categoryID)
	if err != nil {
		return nil, err
	}
	cid, err := budget.NewID(categoryID)
	if err != nil {
		return nil, err
	}

	// TODO handle duplicate error

	id, err := budget.NewID(uuid.New().String())
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
	targetID string,
	amount int,
) error {
	_, err := uuid.Parse(targetID)
	if err != nil {
		return err
	}
	id, err := budget.NewID(targetID)
	if err != nil {
		return err
	}

	if amount < 0 {
		return ErrNegativeNumber
	}

	return s.repo.UpdateAmount(ctx, id, amount)
}

func (s *Service) ReadBudget(
	ctx context.Context,
	id string,
) (*budget.Budget, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return s.repo.FindByID(ctx, budget.ID(id))
}
