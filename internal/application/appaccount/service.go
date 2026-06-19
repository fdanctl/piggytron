package appaccount

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
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
	userID string,
) (*account.Account, error) {
	aid, err := util.ParseID[account.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appaccount.FindOneByID",
		)
		return nil, err
	}
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", uid),
			fmt.Errorf("failed parsing id '%s': %w", uid, err),
			"appaccount.FindOneByID",
		)
		return nil, err
	}

	goal, err := s.repo.FindByID(ctx, aid)
	if err != nil {
		if errors.Is(err, account.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The account does not exists",
				fmt.Errorf("failed to found account '%s': %w", id, err),
				"appaccount.FindOneByID",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to find account '%s': %w", id, err),
				"appaccount.FindOneByID",
			)
		}
		return nil, err
	}
	if goal.UserID() != uid {
		err = errs.NewAppError(
			errs.KindNotFound,
			"The account does not exists",
			fmt.Errorf("the account does not belong to user '%s': %w", uid, account.ErrNotFound),
			"appaccount.FindOneByID",
		)
		return nil, err
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
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.CreateBank",
		)
		return nil, err
	}

	id, err := util.NewID[account.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appuser.CreateBank",
		)
		return nil, err
	}

	acc, err := account.NewBank(id, uid, name, currency, isSaving)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create bank",
			fmt.Errorf("failed to create bank: %w", err),
			"appuser.CreateBank",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, acc)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindValidation,
				"A bank with the same name already exists",
				fmt.Errorf("failed saving account '%s': %w", acc.Name(), err),
				"appuser.CreateBank",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving bank: %w", err),
				"appuser.CreateBank",
			)
		}
		return nil, err
	}
	return acc, nil
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
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.CreateGoal",
		)
		return nil, err
	}

	cid, err := util.ParseID[account.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appaccount.CreateGoal",
		)
		return nil, err
	}

	id, err := util.NewID[account.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appuser.CreateGoal",
		)
		return nil, err
	}

	acc, err := account.NewGoal(
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
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create goal",
			fmt.Errorf("failed to create goal: %w", err),
			"appuser.CreateGoal",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, acc)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindValidation,
				"A goal with the same name already exists",
				fmt.Errorf("failed saving account '%s': %w", acc.Name(), err),
				"appuser.CreateGoal",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving user: %w", err),
				"appuser.CreateGoal",
			)
		}
		return nil, err
	}
	return acc, nil
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
	aid, err := util.ParseID[account.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	cid, err := util.ParseID[account.ID](categoryID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", categoryID),
			fmt.Errorf("failed parsing id '%s': %w", categoryID, err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appuser.UpdateGoal",
		)
		return nil, err
	}
	defer tx.Rollback()

	atx := postgres.NewAccountRepository(tx)
	goal, err := atx.FindByID(ctx, aid)
	if err != nil {
		if errors.Is(err, account.ErrNotFound) {
			err = errs.NewAppError(
				errs.KindNotFound,
				"The goal does not exists",
				fmt.Errorf("failed to found goal '%s': %w", aid, err),
				"appaccount.UpdateGoal",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed finding goal '%s': %w", aid, err),
				"appaccount.UpdateGoal",
			)
		}
		return nil, err
	}
	if goal.UserID() != uid {
		err = errs.NewAppError(
			errs.KindNotFound,
			"The goal does not exists",
			fmt.Errorf("failed to found goal '%s': %w", aid, err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	rtx := postgres.NewTransactionRepository(tx)
	tt, err := rtx.FindAllByAccount(ctx, transaction.ID(goal.ID())) // date DESC
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding account '%s' transactions: %w", goal.ID(), err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}
	var minDate *time.Time
	if len(tt) > 0 {
		d := tt[len(tt)-1].Date()
		minDate = &d
	}

	err = goal.ChangeName(name)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			"Name is invalid",
			fmt.Errorf(
				"failed to changing name of goal '%s' to '%s': %w",
				goal.ID(),
				name,
				err,
			),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}
	err = goal.ChangeTargetAmount(targetAmount)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			"Target amount is invalid",
			fmt.Errorf(
				"failed to change target amount of goal '%s' to '%d': %w",
				goal.ID(),
				targetAmount,
				err,
			),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}
	err = goal.ChangeStartDate(startDate, minDate)
	if err != nil {
		msg := fmt.Sprintf("%s is not a valid date", startDate.String())
		if errors.Is(err, account.ErrContributionBeforeStartDate) {
			msg = fmt.Sprintf("Exists a contribution before %s", startDate.Format(time.DateOnly))
		}
		err = errs.NewAppError(
			errs.KindValidation,
			msg,
			fmt.Errorf(
				"failed to change start date of goal '%s' to '%s': %w",
				goal.ID(),
				startDate.String(),
				err,
			),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}
	err = goal.ChangeTargetDate(targetDate)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid date", startDate.String()),
			fmt.Errorf(
				"failed to change target date of goal '%s' to '%s': %w",
				goal.ID(),
				startDate.String(),
				err,
			),
			"appaccount.UpdateGoal",
		)
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
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to update transactions: %w", err),
				"appaccount.UpdateGoal",
			)
			return nil, err
		}
	}

	err = atx.Update(ctx, goal)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to update goal '%s': %w", goal.ID(), err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.UpdateGoal",
		)
		return nil, err
	}

	return goal, nil
}

func (s *Service) FindAllByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.FindByID",
		)
		return nil, err
	}

	accounts, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' accounts: %w", uid, err),
			"appaccount.FindByID",
		)
		return nil, err
	}

	return accounts, nil
}

func (s *Service) FindAllBanksByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.FindAllBanksByUser",
		)
		return nil, err
	}

	accounts, err := s.repo.FindAllBanksByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' banks: %w", uid, err),
			"appaccount.FindAllBanksByUser",
		)
		return nil, err
	}

	return accounts, nil
}

func (s *Service) FindAllGoalsByUser(
	ctx context.Context,
	userID string,
) ([]*account.Account, error) {
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.FindAllGoalsByUser",
		)
		return nil, err
	}

	accounts, err := s.repo.FindAllGoalsByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' goals: %w", uid, err),
			"appaccount.FindAllGoalsByUser",
		)
		return nil, err
	}

	return accounts, nil
}
