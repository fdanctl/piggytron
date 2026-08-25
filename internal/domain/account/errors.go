package account

import "errors"

var (
	ErrInvalidID                   = errors.New("invalid id")
	ErrInvalidName                 = errors.New("invalid name")
	ErrInvalidCurrency             = errors.New("invalid currency")
	ErrNegativeNumber              = errors.New("number can't be negative")
	ErrNegativeBalance             = errors.New("negative balance")
	ErrDuplicate                   = errors.New("duplicate")
	ErrAccountWrongType            = errors.New("incorrect type")
	ErrContributionBeforeStartDate = errors.New("contribution before start date")
	ErrNotFound                    = errors.New("not found")
	ErrInvalidDate                 = errors.New("invalid date")
	ErrStartDateAfterTarget        = errors.New("start date after target date")
	ErrAccountHasBalance           = errors.New("account has balance")
	ErrClosedAccount               = errors.New("account is closed")
	ErrInvalidType                 = errors.New("invalid type")
)
