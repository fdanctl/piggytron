package views

import (
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Goal struct {
	ID           string
	Name         string
	Type         string
	TargetAmount string
	StartDate    time.Time
	TargetDate   time.Time
	Category     string
	Amount       string

	AmountLeft         string
	MonthlyNeeded      string
	MonthsLeft         string
	CompletePercentage float64
}

func NewGoal(
	g query.AccountWithSum,
) Goal {
	monthlyNeeded := "-"
	monthsLeft := "-"
	if g.TargetDate != nil {
		m := g.TargetDate.Month()

		monthsLeft = "exceded"
		mn := *g.TargetAmount - g.Sum

		if time.Until(*g.TargetDate) > 0 {
			ml := int(m - time.Now().Month())
			for ml < 0 {
				ml += 12
			}
			monthsLeft = strconv.Itoa(ml)
			mn = (*g.TargetAmount - g.Sum) / int(ml)
		}

		monthlyNeeded = formatMoney(float64(mn)/100, currency.EUR, language.AmericanEnglish)
	}

	return Goal{
		ID:   g.ID,
		Name: g.Name,
		Type: g.Type,
		TargetAmount: formatMoney(
			float64(*g.TargetAmount)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		StartDate:  g.CreatedAt,
		TargetDate: *g.TargetDate,
		Category:   g.Category.Name,
		Amount: formatMoney(
			float64(g.Sum)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		AmountLeft: formatMoney(
			float64(*g.TargetAmount-g.Sum)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		MonthlyNeeded:      monthlyNeeded,
		MonthsLeft:         monthsLeft,
		CompletePercentage: float64(g.Sum) / float64(*g.TargetAmount) * 100,
	}
}
