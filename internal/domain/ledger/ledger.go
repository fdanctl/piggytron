// Package ledger defines the ledger entry aggregate — income, expense and
// transfer — the authoritative record of every money movement. Its invariants
// mirror the CHECK constraints in schema.
package ledger

import (
	"errors"
	"time"
)

// ID is a ledger entry identifier.
type ID string

// Type is a ledger entry type: income, expense or transfer.
type Type string

const (
	income   Type = "income"
	expense  Type = "expense"
	transfer Type = "transfer"
)

// NewType parses a string into a Type, returning ErrInvalidType otherwise.
func NewType(str string) (Type, error) {
	switch str {
	case "income":
		return income, nil

	case "expense":
		return expense, nil

	case "transfer":
		return transfer, nil

	default:
		return "", ErrInvalidType
	}
}

// Entry is a single money movement: money in (income), money out (expense) or
// between accounts (transfer). Amounts are positive, in cents.
type Entry struct {
	id     ID
	userID ID

	ttype Type

	fromAccountID *ID
	toAccountID   *ID

	incomeCategoryID  *ID
	expenseCategoryID *ID

	amount      int
	description string
	date        time.Time
	createdAt   time.Time
}

// NewIncome builds a validated income entry credited to toAccountID.
func NewIncome(
	id ID,
	userID ID,
	toAccountID ID,
	incomeCategoryID ID,
	amount int,
	description string,
	date time.Time,
) (*Entry, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Entry{
		id:                id,
		userID:            userID,
		ttype:             income,
		fromAccountID:     nil,
		toAccountID:       &toAccountID,
		incomeCategoryID:  &incomeCategoryID,
		expenseCategoryID: nil,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

// NewExpense builds a validated expense entry debited from fromAccountID.
// minRunningBalance is the account balance before this entry; the expense is
// rejected if it would make the balance negative.
func NewExpense(
	id ID,
	userID ID,
	fromAccountID ID,
	expenseCategoryID ID,
	amount int,
	description string,
	date time.Time,
	minRunningBalance int, // from that date onwards
) (*Entry, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if minRunningBalance-amount < 0 {
		return nil, ErrNegativeBalance
	}

	now := time.Now()

	return &Entry{
		id:                id,
		userID:            userID,
		ttype:             expense,
		fromAccountID:     &fromAccountID,
		toAccountID:       nil,
		incomeCategoryID:  nil,
		expenseCategoryID: &expenseCategoryID,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

// NewTransfer builds a validated transfer between two accounts. For transfers
// into a goal or a savings account, the expense category must match the
// destination account's category and type.
func NewTransfer(
	id ID,
	userID ID,
	fromAccountID ID,
	toAccountID ID,
	expenseCategoryID *ID,
	amount int,
	description string,
	date time.Time,

	minRunningBalance int,
	toAccountCategoryID *ID,
	toAccountCategoryType string,
	isToAccSavings bool,
) (*Entry, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if minRunningBalance-amount < 0 {
		return nil, ErrNegativeBalance
	}
	if fromAccountID == toAccountID {
		return nil, ErrSameAccountTransfer
	}
	if toAccountCategoryID != nil && // ie. is a goal
		(expenseCategoryID == nil || *toAccountCategoryID != *expenseCategoryID) {
		return nil, ErrGoalCategory
	}
	if isToAccSavings {
		if expenseCategoryID == nil {
			return nil, ErrNotSavingsCategory
		}
		if toAccountCategoryType != "savings" {
			return nil, ErrNotSavingsCategory
		}
	}

	now := time.Now()

	return &Entry{
		id:                id,
		userID:            userID,
		ttype:             transfer,
		fromAccountID:     &fromAccountID,
		toAccountID:       &toAccountID,
		incomeCategoryID:  nil,
		expenseCategoryID: expenseCategoryID,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

// Rehydrate rebuilds an Entry from persistence without re-running
// validation (the database constraints already guard it).
func Rehydrate(
	id ID,
	userID ID,
	ttype Type,
	fromAccountID *ID,
	toAccountID *ID,
	incomeCategoryID *ID,
	expenseCategoryID *ID,
	amount int,
	description string,
	date time.Time,
	createdAt time.Time,
) *Entry {
	return &Entry{
		id:                id,
		userID:            userID,
		ttype:             ttype,
		fromAccountID:     fromAccountID,
		toAccountID:       toAccountID,
		incomeCategoryID:  incomeCategoryID,
		expenseCategoryID: expenseCategoryID,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         createdAt,
	}
}

// ID returns the entry id.
func (t *Entry) ID() ID {
	return t.id
}

// UserID returns the id of the user who owns the entry.
func (t *Entry) UserID() ID {
	return t.userID
}

// Type returns the entry type (income, expense or transfer).
func (t *Entry) Type() Type {
	return t.ttype
}

// FromAccountID returns the source account id, or nil for income entries.
func (t *Entry) FromAccountID() *ID {
	return t.fromAccountID
}

// ToAccountID returns the destination account id, or nil for expense entries.
func (t *Entry) ToAccountID() *ID {
	return t.toAccountID
}

// IncomeCategoryID returns the income category id, or nil for non-income entries.
func (t *Entry) IncomeCategoryID() *ID {
	return t.incomeCategoryID
}

// ExpenseCategoryID returns the expense category id, or nil for income
// entries; for transfers it is set only when the transfer moves money into a
// goal or a savings account.
func (t *Entry) ExpenseCategoryID() *ID {
	return t.expenseCategoryID
}

// Amount returns the amount in cents (always positive).
func (t *Entry) Amount() int {
	return t.amount
}

// Description returns the entry description.
func (t *Entry) Description() string {
	return t.description
}

// Date returns the date the entry applies to.
func (t *Entry) Date() time.Time {
	return t.date
}

// CreatedAt returns when the entry was created.
func (t *Entry) CreatedAt() time.Time {
	return t.createdAt
}

// CanBeDeleted reports whether the entry can be deleted: the destination
// account balance (if provided) must stay non-negative after the amount is
// removed from it.
func (t *Entry) CanBeDeleted(toAccBalance *int) error {
	if toAccBalance != nil && *toAccBalance-t.Amount() < 0 {
		return ErrNegativeBalance
	}
	return nil
}

// UpdateIncome replaces the mutable fields of an income entry in place.
func (t *Entry) UpdateIncome(
	toAccountID ID,
	incomeCategoryID ID,
	amount int,
	description string,
	date time.Time,
) error {
	if t.ttype != income {
		return errors.New("can't update income, type not income")
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}

	if description == "" {
		return ErrInvalidDescription
	}

	t.toAccountID = &toAccountID
	t.incomeCategoryID = &incomeCategoryID
	t.amount = amount
	t.description = description
	t.date = date

	return nil
}

// ChangeExpenseCategory reassigns the expense category of a transfer.
func (t *Entry) ChangeExpenseCategory(cid ID) error {
	if t.fromAccountID == nil {
		return errors.New("can't update")
	}

	t.expenseCategoryID = &cid
	return nil
}

// UpdateExpense replaces the mutable fields of an expense entry in place.
func (t *Entry) UpdateExpense(
	fromAccountID ID,
	expenseCategoryID ID,
	amount int,
	description string,
	date time.Time,
) error {
	if t.ttype != expense {
		return errors.New("can't update expense, type not expense")
	}
	if description == "" {
		return ErrInvalidDescription
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}

	t.fromAccountID = &fromAccountID
	t.expenseCategoryID = &expenseCategoryID
	t.amount = amount
	t.description = description
	t.date = date

	return nil
}

// UpdateTransfer replaces the mutable fields of a transfer entry in place.
// It does not verify the destination account; use
// UpdateTransferToAccountAndCategory for that.
func (t *Entry) UpdateTransfer(
	fromAccountID ID,
	toAccountID ID,
	amount int,
	description string,
	date time.Time,
) error {
	if t.ttype != transfer {
		return errors.New("can't update transfer, type not transfer")
	}
	if description == "" {
		return ErrInvalidDescription
	}
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if fromAccountID == toAccountID {
		return ErrSameAccountTransfer
	}

	t.toAccountID = &toAccountID
	t.fromAccountID = &fromAccountID
	t.amount = amount
	t.description = description
	t.date = date

	return nil
}

// UpdateTransferToAccountAndCategory reassigns the destination account (and
// optionally the expense category) of a transfer, re-running the goal/savings
// category checks.
func (t *Entry) UpdateTransferToAccountAndCategory(
	expenseCategoryID *ID,
	toAccountID ID,
	toAccountCategoryID *ID,
	toAccountCategoryType string,
	isToAccSavings bool,
) error {
	if t.ttype != transfer {
		return errors.New("can't update transfer, type not transfer")
	}
	if toAccountID == *t.fromAccountID {
		return ErrSameAccountTransfer
	}
	if toAccountCategoryID != nil && // ie. is a goal
		(expenseCategoryID == nil || *toAccountCategoryID != *expenseCategoryID) {
		return ErrGoalCategory
	}
	if isToAccSavings {
		if expenseCategoryID == nil {
			return ErrNotSavingsCategory
		}
		if toAccountCategoryType != "savings" {
			return ErrNotSavingsCategory
		}
	}

	t.toAccountID = &toAccountID
	t.expenseCategoryID = expenseCategoryID

	return nil
}
