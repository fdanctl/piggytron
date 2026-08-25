// Package appaccount implements the account use cases: creating bank accounts
// and goals, updating goals, and listing accounts by user.
package appaccount

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

// Service implements the account use cases. It composes the account
// repository with a *sql.DB so multi-aggregate operations (e.g., UpdateGoal)
// can run inside a single database transaction.
type Service struct {
	repo account.Repository
	db   *sql.DB
}

// NewService wires the account service to its repository and database.
func NewService(repo account.Repository, db *sql.DB) *Service {
	return &Service{repo: repo, db: db}
}

// FindOneByID returns an account owned by userID, mapping not-found and
// ownership mismatches to KindNotFound.
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

// CreateBank builds and persists a new bank account for the user, rejecting
// duplicate names.
func (s *Service) CreateBank(
	ctx context.Context,
	userID string,
	name string,
	currency string,
	btype string,
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
			"appaccount.CreateBank",
		)
		return nil, err
	}

	accType, err := account.NewType(btype)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid type", btype),
			fmt.Errorf("%s is not a valid type: %w", btype, err),
			"appaccount.CreateBank",
		)
	}

	acc, err := account.NewBank(id, uid, name, currency, accType)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create bank",
			fmt.Errorf("failed to create bank: %w", err),
			"appaccount.CreateBank",
		)
		return nil, err
	}

	err = s.repo.Create(ctx, acc)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf("A %s account with the same name already exists", btype),
				fmt.Errorf("failed saving account '%s': %w", acc.Name(), err),
				"appaccount.CreateBank",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving bank: %w", err),
				"appaccount.CreateBank",
			)
		}
		return nil, err
	}
	return acc, nil
}

// CreateGoal builds and persists a new goal account for the user, rejecting
// duplicate names.
// TODO check if the category is archived
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
			"appaccount.CreateGoal",
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
			"appaccount.CreateGoal",
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
				"appaccount.CreateGoal",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving user: %w", err),
				"appaccount.CreateGoal",
			)
		}
		return nil, err
	}
	return acc, nil
}

// UpdateGoal updates a goal's fields in a single transaction, also
// changes the categories of its funding ledger entries when the
// goal category changes.
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
			"appaccount.UpdateGoal",
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

	rtx := postgres.NewLedgerRepository(tx)
	tt, err := rtx.FindAllByAccount(ctx, ledger.ID(goal.ID())) // date DESC
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding account '%s' ledger entry: %w", goal.ID(), err),
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
			if t.ToAccountID() != nil && *t.ToAccountID() == ledger.ID(goal.ID()) {
				t.ChangeExpenseCategory(ledger.ID(cid))
			}
		}
		err := rtx.UpdateMany(ctx, tt)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to update ledger entry: %w", err),
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

// FindAllByUser returns all accounts of the user.
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

// FindAllBanksByUser returns all bank accounts of the user.
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

// FindAllGoalsByUser returns all goal accounts of the user.
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

func (s *Service) Delete(
	ctx context.Context,
	id string,
) error {
	aid, err := util.ParseID[account.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appledger.Delete",
		)
		return err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appaccount.Delete",
		)
		return err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	hasData := true
	_, err = qtx.GetAccountFirstEntryDate(ctx, id)
	if err != nil {
		if errors.Is(err, account.ErrNotFound) {
			hasData = false
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed finding ledger entries for '%s' account: %w", id, err),
				"appaccount.Delete",
			)
			return err
		}
	}

	if hasData {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Can't delete an account with historical data",
			errors.New("can't delete an account with historical data"),
			"appaccount.Delete",
		)
		return err
	}

	atx := postgres.NewAccountRepository(tx)
	err = atx.Delete(ctx, aid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to delete '%s' account: %w", id, err),
			"appaccount.Delete",
		)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.Delete",
		)
		return err
	}

	return nil
}

func (s *Service) UpdateBankName(
	ctx context.Context,
	userID string,
	id string,
	name string,
	currency string,
	btype string,
) error {
	uid, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.UpdateBankName",
		)
		return err
	}

	aid, err := util.ParseID[account.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.UpdateBankName",
		)
		return err
	}

	accType, err := account.NewType(btype)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid type", btype),
			fmt.Errorf("%s is not a valid type: %w", btype, err),
			"appaccount.Update",
		)
	}

	acc, err := account.NewBank(aid, uid, name, currency, accType)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create bank",
			fmt.Errorf("failed to create bank: %w", err),
			"appaccount.UpdateBankName",
		)
		return err
	}

	// TODO make it just update name
	err = s.repo.Update(ctx, acc)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				fmt.Sprintf("A %s account with the same name already exists", btype),
				fmt.Errorf("failed saving account '%s': %w", acc.Name(), err),
				"appaccount.UpdateBankName",
			)
		} else {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed updating bank: %w", err),
				"appaccount.UpdateBankName",
			)
		}
		return err
	}
	return nil
}

func (s *Service) CloseAccount(
	ctx context.Context,
	id string,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appaccount.CloseAccount",
		)
		return err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	dto, err := qtx.FindWithSum(ctx, id)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding ledger entries for '%s' account: %w", id, err),
			"appaccount.CloseAccount",
		)
		return err
	}

	cid := &dto.Category.ID
	if *cid == util.ZeroUUID {
		cid = nil
	}

	acc := account.Rehydrate(
		account.ID(dto.ID),
		account.ID(dto.UserID),
		account.AccountType(dto.Type),
		dto.Name,
		account.AccountStatus(dto.Status),
		dto.TargetAmount,
		dto.StartDate,
		dto.TargetDate,
		(*account.ID)(cid),
		dto.CompletedAt,
		dto.CancelledAt,
		dto.Currency,
		dto.ClosedAt,
		dto.CreatedAt,
		dto.UpdatedAt,
	)

	err = acc.CloseAccount(dto.Sum)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Make the account balance 0 before closing it",
			err,
			"appaccount.CloseAccount",
		)
		return err
	}

	atx := postgres.NewAccountRepository(tx)
	err = atx.UpdateStatus(ctx, acc)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to close '%s' account: %w", id, err),
			"appaccount.CloseAccount",
		)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CloseAccount",
		)
		return err
	}

	return nil
}
