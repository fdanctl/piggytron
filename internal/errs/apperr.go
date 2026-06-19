package errs

type ErrorKind int

const (
	KindValidation   ErrorKind = iota // 422 – form re-render
	KindBusinessRule                  // 422 – toast
	KindNotFound                      // 404 – toast
	KindConflict                      // 409 – toast (e.g., duplicate)
	KindBadRequest                    // 400 - toast
	KindUnauthorized                  // 401 – redirect to login
	KindInternal                      // 500 – toast, generic message
)

type AppError struct {
	Kind      ErrorKind
	Message   string // user-facing message (for toast)
	Err       error  // underlying error (for logging, wrapping)
	Operation string // e.g., "CreateTransaction"
}

func NewAppError(k ErrorKind, msg string, err error, op string) *AppError {
	return &AppError{
		Kind:      k,
		Message:   msg,
		Err:       err,
		Operation: op,
	}
}

func NewGenericBadRequestAppError(err error, op string) *AppError {
	return &AppError{
		Kind:      KindBadRequest,
		Message:   "Bad request",
		Err:       err,
		Operation: op,
	}
}

func NewInternalAppError(err error, op string) *AppError {
	return &AppError{
		Kind:      KindInternal,
		Message:   "Something went wrong",
		Err:       err,
		Operation: op,
	}
}

func (e *AppError) Error() string { return e.Err.Error() }
func (e *AppError) Unwrap() error { return e.Err }
