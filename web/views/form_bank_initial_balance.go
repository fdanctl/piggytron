package views

import (
	"strconv"
	"strings"
)

type BankInitialBalanceForm struct {
	Form

	TransactionID  string
	InitialBalance string
}

func NewInitialBalanceBankForm() *BankInitialBalanceForm {
	f := BankInitialBalanceForm{}
	f.Initial = true
	return &f
}

func (v *BankInitialBalanceForm) ValidateInitialBalance() (msgs []string) {
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

func (v *BankInitialBalanceForm) InitialBalanceHasError() bool {
	return len(v.ValidateInitialBalance()) > 0
}

func (v *BankInitialBalanceForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateInitialBalance()...)
	return msgs
}
