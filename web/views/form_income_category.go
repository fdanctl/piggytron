package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
)

type IncomeCategoryForm struct {
	Form
	Name string
}

func NewIncomeCategoryForm() *IncomeCategoryForm {
	f := IncomeCategoryForm{}
	f.Initial = true
	return &f
}

func (v *IncomeCategoryForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	if errors.Is(v.CustomError, incomecategory.ErrDuplicate) {
		msgs = append(msgs, "An income category with the same name already exists")
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
