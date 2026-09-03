package views

import (
	"errors"
	"strconv"
	"strings"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"golang.org/x/text/currency"
)

// BankForm is the view-model for the bank account form.
type BankForm struct {
	Form

	Name           string
	Currency       string
	Type           string
	InitialBalance string
}

// NewBankForm returns a blank bank form pre-filled with the EUR currency.
func NewBankForm() *BankForm {
	f := BankForm{
		Currency: currency.EUR.String(),
	}
	f.Initial = true
	f.Type = "checking"
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
		msgs = append(msgs, v.ErrorMsg)
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

func (v *BankForm) ValidateType() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Type != "checking" && v.Type != "savings" {
		msgs = append(msgs, v.Type+" is not a valid type")
	}

	return msgs
}

func (v *BankForm) TypeHasError() bool {
	return len(v.ValidateType()) > 0
}

func (v *BankForm) ValidateInitialBalance() (msgs []string) {
	if v.Initial {
		return
	}

	if v.InitialBalance == "" {
		msgs = append(msgs, "Balance is required")
	}

	str := strings.ReplaceAll(v.InitialBalance, ",", "")
	str = strings.Replace(str, ".", "", 1)

	n, err := strconv.Atoi(str)
	if err != nil {
		return append(msgs, "Not a valid number")
	}

	if n <= 0 {
		msgs = append(msgs, "Balance must greater than 0")
	}
	return msgs
}

func (v *BankForm) InitialBalanceHasError() bool {
	return len(v.ValidateInitialBalance()) > 0
}

func (v *BankForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	msgs = append(msgs, v.ValidateCurrency()...)
	msgs = append(msgs, v.ValidateType()...)
	msgs = append(msgs, v.ValidateInitialBalance()...)
	return msgs
}
