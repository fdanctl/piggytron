package views

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type TransferBudgetForm struct {
	Form
	Amount string
	FromID string
	ToID   string
	Month  string
}

func NewTransferBudgetForm() *TransferBudgetForm {
	f := TransferBudgetForm{}
	f.Initial = true
	return &f
}

func (v *TransferBudgetForm) ValidateAmount() (msgs []string) {
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

	if n <= 0 {
		msgs = append(msgs, "Amount must greater than 0")
	}
	return msgs
}

func (v *TransferBudgetForm) AmountHasError() bool {
	return len(v.ValidateAmount()) > 0
}

func (v *TransferBudgetForm) ValidateFromID() (msgs []string) {
	if v.Initial {
		return
	}

	if v.FromID == "" {
		return append(msgs, "Destination is required")
	}

	if v.FromID != "rta" {
		if _, err := uuid.Parse(v.FromID); err != nil {
			msgs = append(msgs, "Invalid from account")
		}
	}

	return msgs
}

func (v *TransferBudgetForm) FromIDHasError() bool {
	return len(v.ValidateFromID()) > 0
}

func (v *TransferBudgetForm) ValidateToID() (msgs []string) {
	if v.Initial {
		return
	}

	if v.ToID == "" {
		return append(msgs, "Destination is required")
	}

	if v.ToID != "rta" {
		if _, err := uuid.Parse(v.ToID); err != nil {
			msgs = append(msgs, "Invalid to account")
		}
	}

	return msgs
}

func (v *TransferBudgetForm) ToIDHasError() bool {
	return len(v.ValidateToID()) > 0
}

func (v *TransferBudgetForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateAmount()...)
	msgs = append(msgs, v.ValidateFromID()...)
	msgs = append(msgs, v.ValidateToID()...)
	return msgs
}
