// Package appaccount implements the account use cases: creating bank accounts
// and goals, updating goals, and listing accounts by user.
package appaccount

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/account"
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
			"appaccount.FindAllByUser",
		)
		return nil, err
	}

	accounts, err := s.repo.FindAllByUser(ctx, uid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding user '%s' accounts: %w", uid, err),
			"appaccount.FindAllByUser",
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
