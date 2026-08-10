package views

import (
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// Goal is the view-model for a savings goal: current amount, amount left,
// monthly contribution needed and completion percentage.
type Goal struct {
	ID           string
	Name         string
	Type         string
	TargetAmount int
	StartDate    time.Time
	TargetDate   *time.Time
	Category     string
	Amount       string

	AmountLeft         string
	MonthlyNeeded      string
	MonthsLeft         string
	CompletePercentage float64
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

	return Goal{
		ID:           g.ID,
		Name:         g.Name,
		Type:         g.Type,
		TargetAmount: *g.TargetAmount,
		StartDate:    *g.StartDate,
		TargetDate:   g.TargetDate,
		Category:     g.Category.Name,
		Amount: FormatMoney(
			float64(g.Sum)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		AmountLeft: FormatMoney(
			float64(*g.TargetAmount-g.Sum)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		MonthlyNeeded:      monthlyNeeded,
		MonthsLeft:         monthsLeft,
		CompletePercentage: float64(g.Sum) / float64(*g.TargetAmount) * 100,
	}
}
