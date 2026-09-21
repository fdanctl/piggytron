package views

import (
	"strconv"
	"strings"
)

type BudgetHoldForm struct {
	Form
	Month  string
	Amount string
}

func NewBudgetHoldForm() *BudgetHoldForm {
	f := BudgetHoldForm{}
	f.Initial = true
	return &f
}

func (v *BudgetHoldForm) ValidateAmount() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Amount == "" {
		msgs = append(msgs, "Amount is required")
	}

	str := strings.ReplaceAll(v.Amount, ",", "")
	str = strings.Replace(str, ".", "", 1)

	n, err := strconv.Atoi(str)
	if err != nil {
		return append(msgs, "Not a valid number")
	}

	if n < 0 {
		msgs = append(msgs, "Amount must greater than 0")
	}
	return msgs
}

func (v *BudgetHoldForm) AmountHasError() bool {
	return len(v.ValidateAmount()) > 0
}

func (v *BudgetHoldForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateAmount()...)
	return msgs
}
