package budget

import "errors"

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidAmount = errors.New("invalid amount")
)
