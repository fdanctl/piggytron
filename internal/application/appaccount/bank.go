package appaccount

import (
	"context"
	"errors"
	"fmt"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/infrastructure/postgres"
	"github.com/fdanctl/piggytron/internal/util"
)

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
