package appaccount

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/domain/monthlysummary"
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
	initialBalance int,
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

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating account: %w", err),
			"appaccount.CreateBank",
		)
		return nil, err
	}
	defer tx.Rollback()

	atx := postgres.NewAccountRepository(tx)

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

	err = atx.Create(ctx, acc)
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

	if initialBalance > 0 {
		ltx := postgres.NewLedgerRepository(tx)
		mstx := postgres.NewMonthlySummaryRepository(tx)

		entryID, err := util.NewID[ledger.ID]()
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed generating id: %w", err),
				"appaccount.CreateBank",
			)
			return nil, err
		}
		entry, err := ledger.NewBankInitalBalance(
			entryID,
			ledger.ID(uid),
			ledger.ID(id),
			initialBalance,
			"Initial Balance",
			time.Now(),
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create initial balance",
				fmt.Errorf("failed to create interest: %w", err),
				"appaccount.CreateBank",
			)
			return nil, err
		}

		err = ltx.Create(ctx, entry)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving transaction: %w", err),
				"appaccount.CreateBank",
			)
			return nil, err
		}

		// monthly summary
		// toAcc
		tms, err := monthlysummary.New(
			monthlysummary.ID(id),
			monthlysummary.NewMonth(time.Now()),
			initialBalance,
			0,
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create summary",
				fmt.Errorf("failed to create summary: %w", err),
				"appaccount.CreateBank",
			)
			return nil, err
		}

		err = mstx.Save(ctx, tms)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving summary: %w", err),
				"appaccount.CreateBank",
			)
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CreateBank",
		)
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
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
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

func (s *Service) UpdateBankInitial(
	ctx context.Context,
	userID string,
	id string,
	tid string,
	initialBalance int,
) error {
	_, err := util.ParseID[account.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	_, err = util.ParseID[account.ID](id)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", id),
			fmt.Errorf("failed parsing id '%s': %w", id, err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed creating account: %w", err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}
	defer tx.Rollback()

	aqtx := postgres.NewAccountQueryService(tx)
	ltx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)

	var prev *ledger.Entry
	var entryID ledger.ID
	if tid != "" {
		entryID, err = util.ParseID[ledger.ID](tid)
		if err != nil {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s is not a valid id", tid),
				fmt.Errorf("failed parsing id '%s': %w", tid, err),
				"appaccount.UpdateBankInitial",
			)
			return err
		}

		prev, err = ltx.FindByID(ctx, entryID)
		if err != nil {
			if !errors.Is(err, ledger.ErrNotFound) {
				err = errs.NewInternalAppError(
					fmt.Errorf("failed to get first entry date': %w", err),
					"appaccount.UpdateBankInitial",
				)
			}
		}
	} else {
		entryID, err = util.NewID[ledger.ID]()
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed generating id: %w", err),
				"appaccount.UpdateBankInitial",
			)
			return err
		}
	}

	// always the initial-balance date, even if does not exists
	d, err := aqtx.GetAccountFirstEntryDate(ctx, id)
	if err != nil {
		if !errors.Is(err, account.ErrNotFound) {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed to get first entry date': %w", err),
				"appaccount.UpdateBankInitial",
			)
		}
	}

	if d.IsZero() {
		d = time.Now()
	}

	entry, err := ledger.NewBankInitalBalance(
		entryID,
		ledger.ID(userID),
		ledger.ID(id),
		initialBalance,
		"Initial Balance",
		d,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create initial balance",
			fmt.Errorf("failed to create interest: %w", err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	err = ltx.Save(ctx, entry)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to update: %w", err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	// monthly summary
	// toAcc
	delta := initialBalance
	if prev != nil {
		delta -= prev.Amount()
	}
	tms, err := monthlysummary.New(
		monthlysummary.ID(id),
		monthlysummary.NewMonth(d),
		delta,
		0,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	err = mstx.Save(ctx, tms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appaccount.UpdateBankInitial",
		)
		return err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.UpdateBankInitial",
		)
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
		dto.FinalizedAmount,
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
