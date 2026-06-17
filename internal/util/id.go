package util

import (
	"errors"

	"github.com/google/uuid"
)

const ZeroUUID = "00000000-0000-0000-0000-000000000000"

var ErrInvalidID = errors.New("invalid id")

// ParseID validates s as a UUID and returns it typed as T.
func ParseID[T ~string](str string) (T, error) {
	var zero T
	if str == "" {
		return zero, ErrInvalidID
	}
	_, err := uuid.Parse(str)
	if err != nil {
		return zero, ErrInvalidID
	}
	return T(str), nil
}

// NewID generates a new UUID v4 string and returns it as the requested
// domain ID type T.
func NewID[T ~string]() (T, error) {
	return ParseID[T](uuid.New().String())
}
