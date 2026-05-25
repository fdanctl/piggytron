package transaction

import "time"

type ID string

func NewID(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

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

type Transaction struct {
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
) (*Transaction, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Transaction{
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
) (*Transaction, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Transaction{
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
) (*Transaction, error) {
	if description == "" {
		return nil, ErrInvalidDescription
	}
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	now := time.Now()

	return &Transaction{
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
) *Transaction {
	return &Transaction{
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

func (t *Transaction) ID() ID {
	return t.id
}

func (t *Transaction) UserID() ID {
	return t.userID
}

func (t *Transaction) Type() Type {
	return t.ttype
}

func (t *Transaction) FromAccountID() *ID {
	return t.fromAccountID
}

func (t *Transaction) ToAccountID() *ID {
	return t.toAccountID
}

func (t *Transaction) IncomeCategoryID() *ID {
	return t.incomeCategoryID
}

func (t *Transaction) ExpenseCategoryID() *ID {
	return t.expenseCategoryID
}

func (t *Transaction) Amount() int {
	return t.amount
}

func (t *Transaction) Description() string {
	return t.description
}

func (t *Transaction) Date() time.Time {
	return t.date
}

func (t *Transaction) CreatedAt() time.Time {
	return t.createdAt
}
