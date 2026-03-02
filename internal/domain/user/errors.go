package user

import "errors"

var (
	ErrInvalidID       = errors.New("invalid id")
	ErrInvalidName     = errors.New("invalid name")
	ErrInvalidPassword = errors.New("invalid password")
)
