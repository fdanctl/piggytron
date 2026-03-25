package transaction

import "errors"

var (
	ErrInvalidID          = errors.New("invalid id")
	ErrInvalidDescription = errors.New("invalid description")
	ErrInvalidAmount      = errors.New("invalid amount")
	ErrInvalidType        = errors.New("invalid type")
)
