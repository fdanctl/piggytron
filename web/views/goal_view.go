package views

import (
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// GoalStatus marks the status of the account (active or closed).
type GoalStatus string

const (
	ActiveStatus    GoalStatus = "active"
	FundedStatus    GoalStatus = "funded"
	CompletedStatus GoalStatus = "completed"
	CancelledStatus GoalStatus = "cancelled"
)

// Goal is the view-model for a savings goal: current amount, amount left,
// monthly contribution needed and completion percentage.
type Goal struct {
	ID           string
	Name         string
	Status       GoalStatus
	Type         string
	TargetAmount int
	StartDate    time.Time
	TargetDate   *time.Time
	Category     string
	Amount       int

	AmountLeft         int
	MonthlyNeeded      string
	MonthsLeft         string
	CompletePercentage float64
	ClosedAt           *time.Time
}

// NewGoal builds the view-model from the account read-model, deriving the
// monthly contribution needed from the remaining months until the target
// date.
func NewGoal(
	g query.AccountWithSum,
) Goal {
	monthlyNeeded := "-"
	monthsLeft := "-"
	if g.TargetDate != nil {
		monthsLeft = "exceded"
		mn := *g.TargetAmount - g.Sum

		ml := MonthDiff(time.Now(), *g.TargetDate)

		if time.Until(*g.TargetDate) > 0 {
			monthsLeft = strconv.Itoa(ml)
			if ml <= 0 {
				ml++
			}
			mn = (*g.TargetAmount - g.Sum) / ml
		}

		monthlyNeeded = FormatMoney(float64(mn)/100, currency.EUR, language.AmericanEnglish)
	}

	amount := g.Sum
	status := ActiveStatus
	if g.TargetAmount != nil && g.Sum >= *g.TargetAmount {
		status = FundedStatus
	}
	if g.CompletedAt != nil {
		status = CompletedStatus
		amount = *g.FinalizedAmount
	}
	if g.CancelledAt != nil {
		status = CancelledStatus
		amount = *g.FinalizedAmount
	}

	return Goal{
		ID:                 g.ID,
		Name:               g.Name,
		Status:             status,
		Type:               g.Type,
		TargetAmount:       *g.TargetAmount,
		StartDate:          *g.StartDate,
		TargetDate:         g.TargetDate,
		Category:           g.Category.Name,
		Amount:             amount,
		AmountLeft:         *g.TargetAmount - amount,
		MonthlyNeeded:      monthlyNeeded,
		MonthsLeft:         monthsLeft,
		CompletePercentage: float64(amount) / float64(*g.TargetAmount) * 100,
		ClosedAt:           g.ClosedAt,
	}
}
