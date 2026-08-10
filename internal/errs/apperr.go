// Package errs defines the application error taxonomy. AppError wraps a Kind
// (validation, business rule, not found, conflict, bad request, unauthorized,
// internal), a user-facing Message, the underlying error and the operation
// that failed. internal/interface/http/httperror maps kinds to HTTP responses.
package errs

// ErrorKind classifies an application error; httperror maps each kind to an
// HTTP response.
type ErrorKind int

const (
	// KindValidation signals invalid form input; rendered as a 422 form re-render.
	KindValidation ErrorKind = iota
	// KindBusinessRule signals a violated domain rule; 422, shown as a toast.
	KindBusinessRule
	// KindNotFound signals a missing resource; 404, shown as a toast.
	KindNotFound
	// KindConflict signals a uniqueness clash (e.g., duplicate); 409, toast.
	KindConflict
	// KindBadRequest signals a malformed request; 400, shown as a toast.
	KindBadRequest
	// KindUnauthorized signals a missing/expired session; 401, redirect to login.
	KindUnauthorized
	// KindInternal signals an unexpected failure; 500, generic message.
	KindInternal
)

// AppError carries the error kind, a user-facing message, the underlying
// error and the operation that failed.
type AppError struct {
	Kind      ErrorKind
	Message   string // user-facing message (for toast)
	Err       error  // underlying error (for logging, wrapping)
	Operation string // e.g., "CreateTransaction"
}

// NewAppError builds a fully-specified AppError.
func NewAppError(k ErrorKind, msg string, err error, op string) *AppError {
	return &AppError{
		Kind:      k,
		Message:   msg,
		Err:       err,
		Operation: op,
	}
}

// NewGenericBadRequestAppError builds a KindBadRequest AppError with a
// generic "Bad request" message.
func NewGenericBadRequestAppError(err error, op string) *AppError {
	return &AppError{
		Kind:      KindBadRequest,
		Message:   "Bad request",
		Err:       err,
		Operation: op,
	}
}

// NewInternalAppError builds a KindInternal AppError with a generic
// "Something went wrong" message (details stay server-side).
func NewInternalAppError(err error, op string) *AppError {
	return &AppError{
		Kind:      KindInternal,
		Message:   "Something went wrong",
		Err:       err,
		Operation: op,
	}
}

// Error reports the underlying error message (for logs).
func (e *AppError) Error() string { return e.Err.Error() }

// Unwrap exposes the underlying error to errors.Is/As.
func (e *AppError) Unwrap() error { return e.Err }
