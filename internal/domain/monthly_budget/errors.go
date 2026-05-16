package budget

import "errors"

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidAmount = errors.New("invalid amount")
	ErrDuplicate     = errors.New("already exist a budget for this month")
)
