package views

import (
	"errors"

	expensecategory "github.com/fdanctl/piggytron/internal/application/expense_category"
)

type ExpenseCategoryForm struct {
	Initial bool

	Name        string
	Type        int8
	ErrorMsg    string
	CustomError error
}

func NewExpenseCategoryForm() *ExpenseCategoryForm {
	return &ExpenseCategoryForm{
		Initial: true,
	}
}

func (v *ExpenseCategoryForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}
	if errors.Is(v.CustomError, expensecategory.ErrDuplicate) {
		msgs = append(msgs, v.CustomError.Error())
	}
	if len(v.Name) > 30 {
		msgs = append(msgs, "Max length is 30 character")
	}
	return msgs
}

func (v *ExpenseCategoryForm) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *ExpenseCategoryForm) ValidateType() (msgs []string) {
	if v.Initial {
		return
	}
	if v.Type < 1 || v.Type > 3 {
		msgs = append(msgs, "Invalid type")
	}
	return msgs
}

func (v *ExpenseCategoryForm) TypeHasError() bool {
	return len(v.ValidateType()) > 0
}

func (v *ExpenseCategoryForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	return msgs
}
