package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/errs"
)

type Form struct {
	Initial     bool
	ErrorMsg    string
	CustomError error
}

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
