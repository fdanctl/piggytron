package views

import (
	"fmt"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Goal struct {
	Id           account.ID
	Name         string
	Type         account.AccountType
	TargetAmount string
	TargetDate   string
	// Category string
	Amount string

	MonthlyNeeded      string
	MonthsLeft         string
	CompletePercentage float64
}

func NewGoal(
	g *account.Account,
	transactions []*transaction.Transaction,
) Goal {
	var amount int
	for _, t := range transactions {
		if t.ToAccountId() == transaction.ID(g.ID()) {
			amount += t.Amount()
		} else {
			amount -= t.Amount()
		}
	}

	date := "-"
	monthlyNeeded := "-"
	monthsLeft := "-"
	if g.TargetDate() != nil {
		y, m, _ := g.TargetDate().Date()
		date = fmt.Sprintf("%s %s", m, strconv.Itoa(y))

		ml := int(m - time.Now().Month())
		for ml < 0 {
			ml += 12
		}
		monthsLeft = strconv.Itoa(ml)
		mn := (*g.TargetAmount() - amount) / int(ml)
		monthlyNeeded = formatMoney(float64(mn)/100, currency.EUR, language.AmericanEnglish)
	}

	return Goal{
		Id:   g.ID(),
		Name: g.Name(),
		Type: g.Type(),
		TargetAmount: formatMoney(
			float64(*g.TargetAmount())/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		TargetDate: date,
		Amount: formatMoney(
			float64(amount)/100,
			currency.EUR,
			language.AmericanEnglish,
		),
		MonthlyNeeded:      monthlyNeeded,
		MonthsLeft:         monthsLeft,
		CompletePercentage: float64(amount) / float64(*g.TargetAmount()) * 100,
	}
}
