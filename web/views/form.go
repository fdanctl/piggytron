package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/errs"
)

// Form is the base view-model for every form: it carries whether the form
// has been submitted at least once, plus the error message to display.
type Form struct {
	Initial     bool
	ErrorMsg    string
	CustomError error
}

// SetError extracts the user-facing message from an error: AppErrors expose
// their Message, everything else falls back to a generic text.
func (f *Form) SetError(err error) {
	var apperr *errs.AppError
	if errors.As(err, &apperr) {
		f.ErrorMsg = apperr.Message
		f.CustomError = apperr.Err
	} else {
		f.ErrorMsg = "Something went wrong"
		f.CustomError = err
	}
}
