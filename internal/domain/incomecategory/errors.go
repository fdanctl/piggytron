package incomecategory

import "errors"

var (
	ErrInvalidID   = errors.New("invalid id")
	ErrInvalidName = errors.New("invalid name")
	ErrDuplicate   = errors.New("duplicate category name")
	ErrNotFound    = errors.New("not found")
)
