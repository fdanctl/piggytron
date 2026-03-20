package transaction

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidType        = errors.New("invalid type")
	ErrNotNumber          = errors.New("not a number")
	ErrNegativeNumber     = errors.New("number can't be negative")
	ErrMinMax             = errors.New("max in lower than min")
)
