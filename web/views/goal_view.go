package views

import (
	"fmt"
	"strconv"

	"github.com/fdanctl/piggytron/internal/domain/account"
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

	MonthlyNeeded int
	MonthsLeft    int
}

// TODO pass transactions and calculate amount
func NewGoal(
	g *account.Account,
	amount int,
) Goal {
	date := "-"
	if g.TargetDate() != nil {
		y, m, _ := g.TargetDate().Date()
		date = fmt.Sprintf("%s %s", m, strconv.Itoa(y))
	}

	// MonthlyNeeded = amount left / months left
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
		Amount:     formatMoney(float64(amount)/100, currency.EUR, language.AmericanEnglish),
	}
}
