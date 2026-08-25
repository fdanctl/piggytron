// Package account defines the bank and goal aggregates. Banks and goals share
// the Account shape and are discriminated by AccountType; the invariants are
// enforced by New/Rehydrate and mirrored by the CHECK constraints in schema.
package account

import (
	"errors"
	"time"
)

// ID is an account identifier.
type ID string

// AccountType discriminates a bank (checking or savings) from a goal account.
type AccountType string

const (
	CheckingType AccountType = "checking"
	SavingsType  AccountType = "savings"
	GoalType     AccountType = "goal"
)

// AccountStatus marks the status of the account (active or closed).
type AccountStatus string

const (
	ActiveStatus AccountStatus = "active"
	ClosedStatus AccountStatus = "closed"
)

// NewType parses a string into a Type, returning ErrInvalidType otherwise.
func NewType(str string) (AccountType, error) {
	switch str {
	case "checking":
		return CheckingType, nil

	case "savings":
		return SavingsType, nil

	case "goal":
		return GoalType, nil

	default:
		return "", ErrInvalidType
	}
}

// Account is a bank (checking or savings) or a goal. Goals carry target amounts and dates
// plus the expense category used to fund them.
type Account struct {
	id       ID
	userID   ID
	aType    AccountType
	name     string
	status   AccountStatus
	currency string
	// goal-specific
	targetAmount *int
	startDate    *time.Time
	targetDate   *time.Time
	categoryID   *ID
	completedAt  *time.Time
	cancelledAt  *time.Time

	closedAt  *time.Time
	createdAt time.Time
	updatedAt time.Time
}

// NewBank builds a validated bank account.
func NewBank(
	id, userID ID,
	name string,
	currency string,
	btype AccountType,
) (*Account, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if currency == "" || len(currency) > 10 {
		return nil, ErrInvalidCurrency
	}

	now := time.Now()

	return &Account{
		id:        id,
		userID:    userID,
		name:      name,
		aType:     btype,
		status:    ActiveStatus,
		currency:  currency,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// NewGoal builds a validated goal account: target amount must be positive and
// the start date must not be after the target date.
func NewGoal(
	id, userID ID,
	name string,
	currency string,
	targetAmount int,
	startDate time.Time,
	targetDate *time.Time,
	categoryID ID,
) (*Account, error) {
	if name == "" || len(name) > 50 {
		return nil, ErrInvalidName
	}
	if currency == "" || len(currency) > 10 {
		return nil, ErrInvalidCurrency
	}
	if targetAmount <= 0 {
		return nil, ErrNegativeNumber
	}
	if targetDate != nil && startDate.Compare(*targetDate) > 0 {
		return nil, ErrStartDateAfterTarget
	}

	now := time.Now()

	return &Account{
		id:           id,
		userID:       userID,
		name:         name,
		aType:        GoalType,
		status:       ActiveStatus,
		targetAmount: &targetAmount,
		startDate:    &startDate,
		targetDate:   targetDate,
		categoryID:   &categoryID,
		currency:     currency,
		createdAt:    now,
		updatedAt:    now,
	}, nil
}

// Rehydrate rebuilds an Account from persistence without re-running
// validation (the database constraints already guard it).
func Rehydrate(
	id, userID ID,
	aType AccountType,
	name string,
	status AccountStatus,
	targetAmount *int,
	startDate *time.Time,
	targetDate *time.Time,
	categoryID *ID,
	completedAt *time.Time,
	cancelledAt *time.Time,
	currency string,
	closedAt *time.Time,
	createdAt, updatedAt time.Time,
) *Account {
	return &Account{
		id:           id,
		userID:       userID,
		aType:        aType,
		name:         name,
		status:       status,
		targetAmount: targetAmount,
		startDate:    startDate,
		targetDate:   targetDate,
		categoryID:   categoryID,
		completedAt:  completedAt,
		cancelledAt:  cancelledAt,
		currency:     currency,
		closedAt:     closedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ID returns the account id.
func (b *Account) ID() ID {
	return b.id
}

// UserID returns the id of the user who owns the account.
func (b *Account) UserID() ID {
	return b.userID
}

// Name returns the account name.
func (b *Account) Name() string {
	return b.name
}

// Type returns the account type (bank or goal).
func (b *Account) Type() AccountType {
	return b.aType
}

// Status returns the account status (active or closed).
func (b *Account) Status() AccountStatus {
	return b.status
}

// Currency returns the account currency code.
func (b *Account) Currency() string {
	return b.currency
}

// TargetAmount returns the goal target in cents, or nil for banks.
func (b *Account) TargetAmount() *int {
	return b.targetAmount
}

// StartDate returns the goal start date, or nil for banks.
func (b *Account) StartDate() *time.Time {
	return b.startDate
}

// TargetDate returns the goal target date, or nil for banks (unlimited goals).
func (b *Account) TargetDate() *time.Time {
	return b.targetDate
}

// CategoryID returns the goal funding category, or nil for banks.
func (b *Account) CategoryID() *ID {
	return b.categoryID
}

// CompletedAt returns when the goal was completed, or nil if not or if not a goal.
func (b *Account) CompletedAt() *time.Time {
	return b.completedAt
}

// CancelledAt eturns when the goal was cancelled, or nil if not or if not a goal.
func (b *Account) CancelledAt() *time.Time {
	return b.cancelledAt
}

// ClosedAt returns when the account was closed, or nil if not closed.
func (b *Account) ClosedAt() *time.Time {
	return b.closedAt
}

// CreatedAt returns when the account was created.
func (b *Account) CreatedAt() time.Time {
	return b.createdAt
}

// UpdatedAt returns when the account was last updated.
func (b *Account) UpdatedAt() time.Time {
	return b.updatedAt
}

// CanReceiveIncome reports whether the account may be the destination of an
// income entry: goals and savings accounts are funded only by transfers.
func (b *Account) CanReceiveIncome() error {
	if b.IsClosed() {
		return ErrClosedAccount
	}
	if b.aType == GoalType {
		return errors.New("goals can't receive from outside")
	}
	if b.aType == SavingsType {
		return errors.New("savings accounts can't receive from outside")
	}
	return nil
}

// CanMakeExpense reports whether the account may be the source of an expense
// entry: goals and savings accounts are moved only by transfers.
func (b *Account) CanMakeExpense() error {
	if b.IsClosed() {
		return ErrClosedAccount
	}
	if b.aType == GoalType {
		return errors.New("goals can't make expenses")
	}
	if b.aType == SavingsType {
		return errors.New("savings accounts can't make expenses")
	}
	return nil
}

func (b *Account) IsClosed() bool {
	return b.status == ClosedStatus
}

// Update methods: each change re-runs its invariant and bumps UpdatedAt.

// ChangeName updates the account name, requiring a non-empty value up to
// 50 characters long.
func (b *Account) ChangeName(name string) error {
	if name == "" || len(name) > 50 {
		return ErrInvalidName
	}
	b.name = name
	b.updatedAt = time.Now()
	return nil
}

// ChangeTargetAmount updates the goal target, requiring a positive value.
func (b *Account) ChangeTargetAmount(amount int) error {
	if b.aType != GoalType {
		return ErrAccountWrongType
	}
	if amount <= 0 {
		return ErrNegativeNumber
	}
	b.targetAmount = &amount
	b.updatedAt = time.Now()
	return nil
}

// ChangeStartDate moves the goal start date; notLaterThan is the earliest
// contribution date already recorded, which must not precede the new start.
func (b *Account) ChangeStartDate(date time.Time, notLaterThan *time.Time) error {
	if b.aType != GoalType {
		return ErrAccountWrongType
	}
	if notLaterThan != nil && date.Compare(*notLaterThan) == 1 {
		return ErrContributionBeforeStartDate
	}
	if b.TargetDate() != nil && date.Compare(*b.TargetDate()) > 0 {
		return ErrStartDateAfterTarget
	}
	b.startDate = &date
	b.updatedAt = time.Now()
	return nil
}

// ChangeTargetDate updates the goal deadline; nil keeps the goal open-ended.
func (b *Account) ChangeTargetDate(date *time.Time) error {
	if b.aType != GoalType {
		return ErrAccountWrongType
	}
	// targetDate < startDate
	if date != nil && b.StartDate().Compare(*date) > 0 {
		return ErrStartDateAfterTarget
	}
	b.targetDate = date
	b.updatedAt = time.Now()
	return nil
}

// ChangeCategory reassigns the expense category used to fund the goal.
func (b *Account) ChangeCategory(cid ID) error {
	if b.aType != GoalType {
		return ErrAccountWrongType
	}
	b.categoryID = &cid
	b.updatedAt = time.Now()
	return nil
}

// CloseAccount verifies if there is balance and
// closes the account
func (b *Account) CloseAccount(accBalance int) error {
	if accBalance > 0 {
		return ErrAccountHasBalance
	}
	now := time.Now()
	b.status = ClosedStatus
	b.closedAt = &now
	b.updatedAt = now
	return nil
}

func (b *Account) CompleteGoal() error {
	now := time.Now()
	b.status = ClosedStatus
	b.completedAt = &now
	b.closedAt = &now
	b.updatedAt = now
	return nil
}
