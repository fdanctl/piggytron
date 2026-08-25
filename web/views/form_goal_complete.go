package views

import (
	"strconv"
	"strings"
	"time"
)

type GoalCompleteForm struct {
	Form
	Balance      string
	AmountUsed   string
	Date         string
	RemainingDst string
}

func NewGoalCompleteForm(balance int) *GoalCompleteForm {
	f := GoalCompleteForm{
		Date:    time.Now().Format("02/01/2006"),
		Balance: strconv.Itoa(balance),
	}
	f.Initial = true
	return &f
}

func (v *GoalCompleteForm) ValidateAmountUsed() (msgs []string) {
	if v.Initial {
		return
	}

	if v.AmountUsed == "" {
		msgs = append(msgs, "Amount is required")
	}

	str := strings.ReplaceAll(v.AmountUsed, ",", "")
	str = strings.Replace(str, ".", "", 1)

	used, err := strconv.Atoi(str)
	if err != nil {
		return append(msgs, "Not a valid number")
	}

	if used <= 0 {
		msgs = append(msgs, "If not amount is was used, consider cancel the goal instead")
	}

	bal, err := strconv.Atoi(v.Balance)
	if err != nil {
		return
	}

	if bal-used < 0 {
		msgs = append(msgs, "Can't used more than the balance")
	}

	return msgs
}

func (v *GoalCompleteForm) AmountUsedHasError() bool {
	return len(v.ValidateAmountUsed()) > 0
}

func (v *GoalCompleteForm) ValidateDate() (msgs []string) {
	if v.Initial {
		return
	}

	if v.Date == "" {
		msgs = append(msgs, "Date is required")
	}

	_, err := time.Parse("02/01/2006", v.Date)
	if err != nil {
		return append(msgs, "Invalid date")
	}

	return msgs
}

func (v *GoalCompleteForm) DateHasError() bool {
	return len(v.ValidateDate()) > 0
}

func (v *GoalCompleteForm) ValidateRemainingDst() (msgs []string) {
	if v.Initial {
		return
	}

	bal, err := strconv.Atoi(v.Balance)
	if err != nil {
		return
	}

	str := strings.ReplaceAll(v.AmountUsed, ",", "")
	str = strings.Replace(str, ".", "", 1)

	used, err := strconv.Atoi(str)
	if err != nil {
		return
	}

	if bal-used > 0 && v.RemainingDst == "" {
		msgs = append(msgs, "Destination for the remaining is required")
	}

	return msgs
}

func (v *GoalCompleteForm) RemainingDstHasError() bool {
	return len(v.ValidateRemainingDst()) > 0
}

func (v *GoalCompleteForm) Validate() (msgs []string) {
	if v.Initial {
		return
	}
	msgs = append(msgs, v.ValidateAmountUsed()...)
	msgs = append(msgs, v.ValidateDate()...)
	msgs = append(msgs, v.ValidateRemainingDst()...)
	return msgs
}
