package account

import "errors"

var (
	ErrInvalidID                   = errors.New("invalid id")
	ErrInvalidName                 = errors.New("invalid name")
	ErrInvalidCurrency             = errors.New("invalid currency")
	ErrNegativeNumber              = errors.New("number can't be negative")
	ErrDuplicate                   = errors.New("duplicate")
	ErrAccountWrongType            = errors.New("incorrect type")
	ErrContributionBeforeStartDate = errors.New("contribution before start date")
	ErrNotFound                    = errors.New("not found")
	ErrStartDateAfterTarget        = errors.New("start date after target date")
)
