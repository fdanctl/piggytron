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

func (s *Service) CompleteGoal(
	ctx context.Context,
	aid string,
	userID string,
	amount int,
	date time.Time,
	remainingAccID string,
) (*account.Account, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](aid)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", aid),
			fmt.Errorf("failed parsing id '%s': %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	goalAcc, err := qtx.FindWithSum(ctx, aid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding ledger entries for '%s' account: %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	if goalAcc.UserID != userID {
		return nil, account.ErrNotFound
	}

	if goalAcc.Sum-amount < 0 {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Amount used can't be higher than balance",
			fmt.Errorf("amount used can't be higher than balance: %w", account.ErrNegativeBalance),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	lastEntryDate, err := qtx.GetAccountLastEntryDate(ctx, aid)
	if date.Before(lastEntryDate) {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Contributions found after the complete date",
			fmt.Errorf("contributions found after the complete date: %w", account.ErrInvalidDate),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	// remaining transfer
	if goalAcc.Sum-amount > 0 {
		if remainingAccID == "" {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination account is necessary for the remaining",
				fmt.Errorf(
					"empty account id: %w",
					account.ErrInvalidID,
				),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}
		remainingTarget, err := qtx.FindWithSum(ctx, remainingAccID)
		if err != nil {
			if errors.Is(err, account.ErrNotFound) {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					"Remaining destination account not found",
					fmt.Errorf(
						"remaining destination account not found: %w",
						err,
					),
					"appaccount.CompleteGoal",
				)
			} else {
				err = errs.NewInternalAppError(err, "appaccount.CompleteGoal")
			}
			return nil, err
		}

		if remainingTarget.UserID != userID {
			return nil, account.ErrNotFound
		}

		if remainingTarget.Type != string(account.CheckingType) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining amount can only be transfer to a checking account",
				fmt.Errorf(
					"remaining amount can only be transfer to a checking account: %w",
					account.ErrNegativeBalance,
				),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		if remainingTarget.Status == string(account.ClosedStatus) {
			return nil, errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination is closed",
				account.ErrClosedAccount,
				"appaccount.CompleteGoal",
			)
		}

		tid, err := util.NewID[ledger.ID]()
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed generating id: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		toAccID, err := util.ParseID[ledger.ID](remainingTarget.ID)
		if err != nil {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s is not a valid id", remainingTarget.ID),
				fmt.Errorf("failed parsing id '%s': %w", remainingTarget.ID, err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		t, err := ledger.NewTransfer(
			tid,
			uid,
			fromAccID,
			toAccID,
			nil,
			goalAcc.Sum-amount,
			fmt.Sprintf("Remaining from completing %s goal", goalAcc.Name),
			date,
			goalAcc.Sum,
			nil,
			"",
			remainingTarget.Type == string(account.SavingsType),
		)
		if err != nil {
			msg := "Failed to create transfer"
			if errors.Is(err, ledger.ErrNegativeBalance) {
				msg = fmt.Sprintf(
					"%s becomes negative on %s",
					goalAcc.Name,
					date,
				)
			}
			err = errs.NewAppError(
				errs.KindBusinessRule,
				msg,
				fmt.Errorf("failed to create transfer: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		err = rtx.Create(ctx, t)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving transaction: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		// monthly summary
		// toAcc
		tms, err := monthlysummary.New(
			monthlysummary.ID(remainingTarget.ID),
			monthlysummary.NewMonth(date),
			goalAcc.Sum-amount,
			0,
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create summary",
				fmt.Errorf("failed to create summary: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}

		err = mstx.Save(ctx, tms)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving summary: %w", err),
				"appaccount.CompleteGoal",
			)
			return nil, err
		}
	}

	cid := &goalAcc.Category.ID
	if *cid == util.ZeroUUID {
		cid = nil
	}

	acc := account.Rehydrate(
		account.ID(goalAcc.ID),
		account.ID(goalAcc.UserID),
		account.AccountType(goalAcc.Type),
		goalAcc.Name,
		account.AccountStatus(goalAcc.Status),
		goalAcc.TargetAmount,
		goalAcc.StartDate,
		goalAcc.TargetDate,
		(*account.ID)(cid),
		goalAcc.CompletedAt,
		goalAcc.CancelledAt,
		goalAcc.Currency,
		goalAcc.ClosedAt,
		goalAcc.CreatedAt,
		goalAcc.UpdatedAt,
	)

	err = acc.CompleteGoal()
	if err != nil {
		return nil, err
	}

	atx := postgres.NewAccountRepository(tx)
	err = atx.UpdateStatus(ctx, acc)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to complete '%s' account: %w", aid, err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	gfid, err := util.NewID[ledger.ID]()
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed generating id: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	fulfillmentEntry, err := ledger.NewGoalFulfillment(
		gfid,
		uid,
		fromAccID,
		acc.Name(),
		ledger.ID(*acc.CategoryID()),
		amount,
		date,
	)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create transfer",
			fmt.Errorf("failed to create goal fulfillment: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	err = rtx.Create(ctx, fulfillmentEntry)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving transaction: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	// monthly summary
	// fromAcc
	fms, err := monthlysummary.New(
		monthlysummary.ID(aid),
		monthlysummary.NewMonth(date),
		0,
		goalAcc.Sum,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	err = mstx.Save(ctx, fms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CompleteGoal",
		)
		return nil, err
	}

	return acc, nil
}

func (s *Service) CancelGoal(
	ctx context.Context,
	aid string,
	userID string,
	destinationAccID string,
) (*account.Account, error) {
	uid, err := util.ParseID[ledger.ID](userID)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", userID),
			fmt.Errorf("failed parsing id '%s': %w", userID, err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	fromAccID, err := util.ParseID[ledger.ID](aid)
	if err != nil {
		err = errs.NewAppError(
			errs.KindValidation,
			fmt.Sprintf("%s is not a valid id", aid),
			fmt.Errorf("failed parsing id '%s': %w", aid, err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed updating goal: %w", err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}
	defer tx.Rollback()

	qtx := postgres.NewAccountQueryService(tx)
	rtx := postgres.NewLedgerRepository(tx)
	mstx := postgres.NewMonthlySummaryRepository(tx)
	goalAcc, err := qtx.FindWithSum(ctx, aid)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed finding ledger entries for '%s' account: %w", aid, err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	if goalAcc.UserID != userID {
		return nil, account.ErrNotFound
	}

	// transfer
	if goalAcc.Sum > 0 {
		if destinationAccID == "" {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination account is necessary for the remaining",
				fmt.Errorf(
					"empty account id: %w",
					account.ErrInvalidID,
				),
				"appaccount.CancelGoal",
			)
			return nil, err
		}
		destinationAcc, err := qtx.FindWithSum(ctx, destinationAccID)
		if err != nil {
			if errors.Is(err, account.ErrNotFound) {
				err = errs.NewAppError(
					errs.KindBusinessRule,
					"Destination account not found",
					fmt.Errorf(
						"Destination account not found: %w",
						err,
					),
					"appaccount.CancelGoal",
				)
			} else {
				err = errs.NewInternalAppError(err, "appaccount.CancelGoal")
			}
			return nil, err
		}

		if destinationAcc.UserID != userID {
			return nil, account.ErrNotFound
		}

		if destinationAcc.Type != string(account.CheckingType) {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Amount can only be transfer to a checking account",
				fmt.Errorf(
					"Amount can only be transfer to a checking account: %w",
					account.ErrNegativeBalance,
				),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		if destinationAcc.Status == string(account.ClosedStatus) {
			return nil, errs.NewAppError(
				errs.KindBusinessRule,
				"Remaining destination is closed",
				account.ErrClosedAccount,
				"appaccount.CancelGoal",
			)
		}

		tid, err := util.NewID[ledger.ID]()
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed generating id: %w", err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		toAccID, err := util.ParseID[ledger.ID](destinationAcc.ID)
		if err != nil {
			err = errs.NewAppError(
				errs.KindValidation,
				fmt.Sprintf("%s is not a valid id", destinationAcc.ID),
				fmt.Errorf("failed parsing id '%s': %w", destinationAcc.ID, err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		t, err := ledger.NewTransfer(
			tid,
			uid,
			fromAccID,
			toAccID,
			nil,
			goalAcc.Sum,
			fmt.Sprintf("%s goal cancelation", goalAcc.Name),
			time.Now(),
			goalAcc.Sum,
			nil,
			"",
			destinationAcc.Type == string(account.SavingsType),
		)
		if err != nil {
			msg := "Failed to create transfer"
			if errors.Is(err, ledger.ErrNegativeBalance) {
				msg = fmt.Sprintf(
					"%s becomes negative on %s",
					goalAcc.Name,
					time.Now(),
				)
			}
			err = errs.NewAppError(
				errs.KindBusinessRule,
				msg,
				fmt.Errorf("failed to create transfer: %w", err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		err = rtx.Create(ctx, t)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving transaction: %w", err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		// monthly summary
		// toAcc
		tms, err := monthlysummary.New(
			monthlysummary.ID(destinationAcc.ID),
			monthlysummary.NewMonth(time.Now()),
			goalAcc.Sum,
			0,
		)
		if err != nil {
			err = errs.NewAppError(
				errs.KindBusinessRule,
				"Failed to create summary",
				fmt.Errorf("failed to create summary: %w", err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}

		err = mstx.Save(ctx, tms)
		if err != nil {
			err = errs.NewInternalAppError(
				fmt.Errorf("failed saving summary: %w", err),
				"appaccount.CancelGoal",
			)
			return nil, err
		}
	}

	cid := &goalAcc.Category.ID
	if *cid == util.ZeroUUID {
		cid = nil
	}

	acc := account.Rehydrate(
		account.ID(goalAcc.ID),
		account.ID(goalAcc.UserID),
		account.AccountType(goalAcc.Type),
		goalAcc.Name,
		account.AccountStatus(goalAcc.Status),
		goalAcc.TargetAmount,
		goalAcc.StartDate,
		goalAcc.TargetDate,
		(*account.ID)(cid),
		goalAcc.CompletedAt,
		goalAcc.CancelledAt,
		goalAcc.Currency,
		goalAcc.ClosedAt,
		goalAcc.CreatedAt,
		goalAcc.UpdatedAt,
	)

	err = acc.CancelGoal()
	if err != nil {
		return nil, err
	}

	atx := postgres.NewAccountRepository(tx)
	err = atx.UpdateStatus(ctx, acc)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to cancel '%s' account: %w", aid, err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	// monthly summary
	// fromAcc
	fms, err := monthlysummary.New(
		monthlysummary.ID(aid),
		monthlysummary.NewMonth(time.Now()),
		0,
		goalAcc.Sum,
	)
	if err != nil {
		err = errs.NewAppError(
			errs.KindBusinessRule,
			"Failed to create summary",
			fmt.Errorf("failed to create summary: %w", err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	err = mstx.Save(ctx, fms)
	if err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed saving summary: %w", err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		err = errs.NewInternalAppError(
			fmt.Errorf("failed to commit': %w", err),
			"appaccount.CancelGoal",
		)
		return nil, err
	}

	return acc, nil
}
