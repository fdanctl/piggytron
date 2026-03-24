package transaction

import "time"

type ID string

func NewId(str string) (ID, error) {
	if str == "" {
		return "", ErrInvalidID
	}
	return ID(str), nil
}

type Ttype string

const (
	income   Ttype = "income"
	expense  Ttype = "expense"
	transfer Ttype = "transfer"
)

func NewType(str string) (Ttype, error) {
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
	userId ID

	ttype Ttype

	fromAccountId *ID
	toAccountId   *ID

	incomeCategoryId  *ID
	expenseCategoryId *ID

	amount      int
	description string
	date        time.Time
	createdAt   time.Time
}

func NewIncome(
	id ID,
	userId ID,
	toAccountId ID,
	incomeCategoryId ID,
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
		userId:            id,
		ttype:             income,
		fromAccountId:     nil,
		toAccountId:       &toAccountId,
		incomeCategoryId:  &incomeCategoryId,
		expenseCategoryId: nil,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

func NewExpense(
	id ID,
	userId ID,
	fromAccountId ID,
	expenseCategoryId ID,
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
		userId:            id,
		ttype:             expense,
		fromAccountId:     &fromAccountId,
		toAccountId:       nil,
		incomeCategoryId:  nil,
		expenseCategoryId: &expenseCategoryId,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

func NewTransfer(
	id ID,
	userId ID,
	fromAccountId ID,
	toAccountId ID,
	expenseCategoryId *ID,
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
		userId:            id,
		ttype:             transfer,
		fromAccountId:     &fromAccountId,
		toAccountId:       &toAccountId,
		incomeCategoryId:  nil,
		expenseCategoryId: expenseCategoryId,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         now,
	}, nil
}

func Rehydrate(
	id ID,
	userId ID,
	ttype Ttype,
	fromAccountId *ID,
	toAccountId *ID,
	incomeCategoryId *ID,
	expenseCategoryId *ID,
	amount int,
	description string,
	date time.Time,
	createdAt time.Time,
) *Transaction {
	return &Transaction{
		id:                id,
		userId:            id,
		ttype:             ttype,
		fromAccountId:     fromAccountId,
		toAccountId:       toAccountId,
		incomeCategoryId:  incomeCategoryId,
		expenseCategoryId: expenseCategoryId,
		amount:            amount,
		description:       description,
		date:              date,
		createdAt:         createdAt,
	}
}

func (t *Transaction) ID() ID {
	return t.id
}

func (t *Transaction) UserId() ID {
	return t.userId
}

func (t *Transaction) Ttype() Ttype {
	return t.ttype
}

func (t *Transaction) FromAccountId() ID {
	if t.fromAccountId == nil {
		return ID("")
	}
	return *t.fromAccountId
}

func (t *Transaction) ToAccountId() ID {
	if t.toAccountId == nil {
		return ID("")
	}
	return *t.toAccountId
}

func (t *Transaction) IncomeCategoryId() ID {
	if t.incomeCategoryId == nil {
		return ID("")
	}
	return *t.incomeCategoryId
}

func (t *Transaction) ExpenseCategoryId() ID {
	if t.expenseCategoryId == nil {
		return ID("")
	}
	return *t.expenseCategoryId
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
