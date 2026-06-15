package views

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/google/uuid"
	"golang.org/x/text/currency"
)

type GoalForm struct {
	Initial bool

	Name         string
	Currency     string
	TargetAmount string
	StartDate    string
	TargetDate   string
	Category     string

	ErrorMsg    string
	CustomError error
}

func NewGoalForm() *GoalForm {
	return &GoalForm{
		Initial:   true,
		StartDate: time.Now().Format("02/01/2006"),
		Currency:  currency.EUR.String(),
	}
}

func (v *GoalForm) ValidateName() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Name == "" {
		msgs = append(msgs, "Name is required")
	}

	if errors.Is(v.CustomError, account.ErrDuplicate) {
		msgs = append(msgs, v.CustomError.Error())
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
		msgs = append(msgs, v.Currency+" is not a valid currency")
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

	if v.TargetAmount == "" {
		msgs = append(msgs, "Target amount is required")
	}

	str := strings.ReplaceAll(v.TargetAmount, ",", "")
	str = strings.Replace(str, ".", "", 1)

	n, err := strconv.Atoi(str)
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

func (v *GoalForm) ValidateStartDate() (msgs []string) {
	if v.Initial {
		return
	}

	sdate, err := time.Parse("02/01/2006", v.StartDate)
	if err != nil {
		return append(msgs, "Invalid date")
	}

	if v.TargetDate == "" {
		return
	}

	if errors.Is(v.CustomError, account.ErrContributionBeforeStartDate) {
		msgs = append(msgs, v.CustomError.Error())
	}

	tdate, _ := time.Parse("02/01/2006", v.TargetDate)

	duration := tdate.Sub(sdate)
	if duration < 0 {
		msgs = append(msgs, "Start date can't be after target date")
	}

	return msgs
}

func (v *GoalForm) StartDateHasError() bool {
	return len(v.ValidateStartDate()) > 0
}

func (v *GoalForm) ValidateTargetDate() (msgs []string) {
	if v.Initial {
		return
	}

	// optional
	if v.TargetDate == "" {
		return
	}

	date, err := time.Parse("02/01/2006", v.TargetDate)
	if err != nil {
		return append(msgs, "Invalid date")
	}

	duration := time.Until(date)
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
	msgs = append(msgs, v.ValidateStartDate()...)
	msgs = append(msgs, v.ValidateTargetDate()...)
	msgs = append(msgs, v.ValidateCategory()...)
	return msgs
}
