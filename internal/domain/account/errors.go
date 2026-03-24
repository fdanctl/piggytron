package account

import "errors"

var (
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidCurrency = errors.New("invalid currency")
	ErrNegativeNumber  = errors.New("number can't be negative")
)
