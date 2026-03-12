package expensecategory

import "errors"

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidType = errors.New("invalid type")
)
