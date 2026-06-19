package expensecategory

import "errors"

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrInvalidName = errors.New("invalid name")
	ErrInvalidType = errors.New("invalid type")
	ErrDuplicate   = errors.New("duplicate category name")
	ErrNotFound    = errors.New("not found")
)
