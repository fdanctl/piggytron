package views

import (
	"errors"

	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	expensecategory "github.com/fdanctl/piggytron/internal/domain/expense_category"
)

type ExpenseCategoryForm struct {
	Initial bool

	Name        string
	Type        string
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
	if errors.Is(v.CustomError, expensecategoryapp.ErrDuplicate) {
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
	if _, err := expensecategory.NewExpenseType(v.Type); err != nil {
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
	msgs = append(msgs, v.ValidateType()...)
	return msgs
}
