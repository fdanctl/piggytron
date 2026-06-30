package ledger

import (
	"errors"
	"time"
)

type ID string

type Type string

const (
	income   Type = "income"
	expense  Type = "expense"
	transfer Type = "transfer"
)

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

func NewExpense(
	id ID,
	userID ID,
	fromAccountID ID,
	expenseCategoryID ID,
	amount int,
	description string,
	date time.Time,
	sourceBalance int,
) (*Entry, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}
	if sourceBalance-amount < 0 {
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

func NewTransfer(
	id ID,
	userID ID,
	fromAccountID ID,
	toAccountID ID,
	expenseCategoryID *ID,
	amount int,
	description string,
	date time.Time,

	sourceBalance int,
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
	if sourceBalance-amount < 0 {
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

func (t *Entry) ID() ID {
	return t.id
}

func (t *Entry) UserID() ID {
	return t.userID
}

func (t *Entry) Type() Type {
	return t.ttype
}

func (t *Entry) FromAccountID() *ID {
	return t.fromAccountID
}

func (t *Entry) ToAccountID() *ID {
	return t.toAccountID
}

func (t *Entry) IncomeCategoryID() *ID {
	return t.incomeCategoryID
}

func (t *Entry) ExpenseCategoryID() *ID {
	return t.expenseCategoryID
}

func (t *Entry) Amount() int {
	return t.amount
}

func (t *Entry) Description() string {
	return t.description
}

func (t *Entry) Date() time.Time {
	return t.date
}

func (t *Entry) CreatedAt() time.Time {
	return t.createdAt
}

// CanBeDeleted receive the destination account balance or nil if it does not nonexistent
func (t *Entry) CanBeDeleted(toAccBalance *int) error {
	if toAccBalance != nil && *toAccBalance-t.Amount() < 0 {
		return ErrNegativeBalance
	}
	return nil
}

// updates

func (t *Entry) ChangeExpenseCategory(cid ID) error {
	if t.fromAccountID == nil {
		return errors.New("can't update")
	}

	t.expenseCategoryID = &cid
	return nil
}
