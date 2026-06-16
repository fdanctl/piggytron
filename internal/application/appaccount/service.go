package appaccount

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/google/uuid"
)

type Service struct {
	repo account.Repository
	db   *sql.DB
}

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidDate   = errors.New("invalid date")
)

func NewService(repo account.Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

func (s *Service) FindOneByID(
	ctx context.Context,
	id string,
	uid string,
) (*account.Account, error) {
	_, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	accID, err := account.NewID(id)
	if err != nil {
		return nil, err
	}

	goal, err := s.repo.FindByID(ctx, accID)
	if goal.UserID() != account.ID(uid) {
		return nil, errors.New("not found")
	}
	return goal, nil
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	atx := postgres.NewAccountRepository(tx)
	goal, err := atx.FindByID(ctx, account.ID(id))
	if err != nil {
		return nil, err
	}
	if goal.UserID() != uid {
		return nil, errors.New("not found")
	}

	rtx := postgres.NewTransactionRepository(tx)
	tt, err := rtx.FindAllByAccount(ctx, transaction.ID(goal.ID())) // date DESC
	if err != nil {
		return nil, err
	}
	var minDate *time.Time
	if len(tt) > 0 {
		d := tt[len(tt)-1].Date()
		minDate = &d
	}

	err = goal.ChangeName(name)
	if err != nil {
		return nil, err
	}
	err = goal.ChangeTargetAmount(targetAmount)
	if err != nil {
		return nil, err
	}
	err = goal.ChangeStartDate(startDate, minDate)
	if err != nil {
		return nil, err
	}
	err = goal.ChangeTargetDate(targetDate)
	if err != nil {
		return nil, err
	}

	if goal.CategoryID() != nil && *goal.CategoryID() != cid {
		goal.ChangeCategory(cid)

		for _, t := range tt {
			if t.ToAccountID() != nil && *t.ToAccountID() == transaction.ID(goal.ID()) {
				t.ChangeExpenseCategory(transaction.ID(cid))
			}
		}
		err := rtx.UpdateMany(ctx, tt)
		if err != nil {
			return nil, err
		}
	}

	err = atx.Update(ctx, goal)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
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
