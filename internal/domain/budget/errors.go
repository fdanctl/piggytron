package budget

import "errors"

var (
	ErrInvalidID     = errors.New("invalid id")
	ErrInvalidMonth  = errors.New("month should be in YYYY-MM format, ex: 2026-05")
	ErrInvalidAmount = errors.New("invalid amount")
	ErrDuplicate     = errors.New("already exist a budget for this month")
	ErrNotFound      = errors.New("not found")
)
