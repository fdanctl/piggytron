package views

import (
	"errors"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"golang.org/x/text/currency"
)

type BankForm struct {
	Form

	Name      string
	Currency  string
	IsSavings bool
}

func NewBankForm() *BankForm {
	f := BankForm{
		Currency: currency.EUR.String(),
	}
	f.Initial = true
	return &f
}

func (v *BankForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}

	if errors.Is(v.CustomError, account.ErrDuplicate) {
		msgs = append(msgs, "A bank with the same name already exists")
	}

	if len(v.Name) > 50 {
		msgs = append(msgs, "Max length is 50 characters")
	}

	return msgs
}

func (v *BankForm) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *BankForm) ValidateCurrency() (msgs []string) {
	if v.Initial {
		return
	}

	_, err := currency.ParseISO(v.Currency)
	if err != nil {
		msgs = append(msgs, v.Currency+" is not a valid currency")
	}

	return msgs
}

func (v *BankForm) CurrencyHasError() bool {
	return len(v.ValidateCurrency()) > 0
}

func (v *BankForm) ValidateIsSavings() (msgs []string) {
	if v.Initial {
		return
	}

	return msgs
}

func (v *BankForm) IsSavingsHasError() bool {
	return len(v.ValidateIsSavings()) > 0
}

func (v *BankForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	msgs = append(msgs, v.ValidateCurrency()...)
	msgs = append(msgs, v.ValidateIsSavings()...)
	return msgs
}
