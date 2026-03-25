package views

import (
	"strconv"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type GoalForm struct {
	Initial bool

	Name         string
	Currency     string
	TargetAmount string
	TargetDate   string
	Category     string
}

func NewGoalForm() *GoalForm {
	return &GoalForm{
		Initial:  true,
		Currency: currency.EUR.String(),
	}
}

func (v *GoalForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}

	if len(v.Name) > 50 {
		msgs = append(msgs, "Max length is 50 characters")
	}

	return msgs
}

func (v *GoalForm) NameHasError() bool {
	return len(v.ValidateName()) > 0
}

func (v *GoalForm) ValidateCurrency() (msgs []string) {
	if v.Initial {
		return
	}

	_, err := currency.ParseISO(v.Currency)
	if err != nil {
		msgs = append(msgs, v.Category+" is not a valid currency")
	}

	return msgs
}

func (v *GoalForm) CurrencyHasError() bool {
	return len(v.ValidateCurrency()) > 0
}

func (v *GoalForm) ValidateTargetAmount() (msgs []string) {
	if v.Initial {
		return
	}

	n, err := strconv.Atoi(v.TargetAmount)
	if err != nil {
		return append(msgs, "Not a valid number")
	}

	if n < 0 {
		msgs = append(msgs, "Amount can't be negative")
	}
	return msgs
}

func (v *GoalForm) TargetAmountHasError() bool {
	return len(v.ValidateTargetAmount()) > 0
}

func (v *GoalForm) ValidateTargetDate() (msgs []string) {
	if v.Initial {
		return
	}

	// optional
	if v.TargetDate == "" {
		return
	}

	date, err := time.Parse(time.DateOnly, v.TargetDate)
	if err != nil {
		return append(msgs, "Invalid date")
	}

	duration := date.Sub(time.Now())
	if duration < 0 {
		msgs = append(msgs, "It's for yesterday?")
	}

	return msgs
}

func (v *GoalForm) TargetDateHasError() bool {
	return len(v.ValidateTargetDate()) > 0
}

func (v *GoalForm) ValidateCategory() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Category == "" {
		return append(msgs, "Category is required")
	}

	if _, err := uuid.Parse(v.Category); err != nil {
		msgs = append(msgs, "Invalid category")
	}

	return msgs
}

func (v *GoalForm) CategoryHasError() bool {
	return len(v.ValidateCategory()) > 0
}

func (v *GoalForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateName()...)
	msgs = append(msgs, v.ValidateCurrency()...)
	msgs = append(msgs, v.ValidateTargetAmount()...)
	msgs = append(msgs, v.ValidateTargetDate()...)
	msgs = append(msgs, v.ValidateCategory()...)
	return msgs
}
