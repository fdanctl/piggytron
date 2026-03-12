package views

import (
	"errors"

	incomecategory "github.com/fdanctl/piggytron/internal/application/income_category"
)

type IncomeCategoryForm struct {
	Initial bool

	Name        string
	ErrorMsg    string
	CustomError error
}

func NewIncomeCategoryForm() *IncomeCategoryForm {
	return &IncomeCategoryForm{
		Initial: true,
	}
}

func (v *IncomeCategoryForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	if errors.Is(v.CustomError, incomecategory.ErrDuplicate) {
		msgs = append(msgs, v.CustomError.Error())
	}
	if len(v.Name) > 30 {
		msgs = append(msgs, "Max length is 30 character")
	}
	return msgs
}

func (v *IncomeCategoryForm) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *IncomeCategoryForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	return msgs
}
