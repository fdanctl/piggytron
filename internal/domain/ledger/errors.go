package ledger

import "errors"

var (
	ErrInvalidID           = errors.New("invalid id")
	ErrInvalidDescription  = errors.New("invalid description")
	ErrInvalidAmount       = errors.New("invalid amount")
	ErrNegativeBalance     = errors.New("negative balance")
	ErrGoalCategory        = errors.New("wrong category for goal")
	ErrNotSavingsCategory  = errors.New("not a saving category")
	ErrInvalidAccount      = errors.New("income/expense not possible for goals nor savings account")
	ErrInvalidType         = errors.New("invalid type")
	ErrNotFound            = errors.New("not found")
	ErrSameAccountTransfer = errors.New("same account transfer")
)
