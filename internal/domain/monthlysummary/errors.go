package monthlysummary

import "errors"

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrInvalidMonth  = errors.New("month should be in YYYY-MM format, ex: 2026-05")
	ErrNotFound      = errors.New("not found")
)
