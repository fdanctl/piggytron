package views

import (
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Transaction struct {
	ID          string
	Description string
	Type        string
	Amount      string
	Date        string
}

func NewTransaction(
	t query.TransactionDTO,
) Transaction {
	f := float64(t.Amount) / 100
	if t.Type == "expense" {
		f *= -1
	}
	return Transaction{
		ID:          t.ID,
		Description: t.Description,
		Type:        t.Type,
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		// TODO convert date to relative date (ex. today, yesterday, ...)
		Date: t.Date.Format(time.DateOnly),
	}
}

func NewAccountTransaction(
	t query.TransactionDTO,
	a *query.AccountWithCategory,
) Transaction {
	f := float64(t.Amount) / 100
	tpe := "income"
	fmt.Println(t.FromAccountID, &a.ID)
	if t.FromAccountID != nil && *t.FromAccountID == a.ID {
		tpe = "expense"
		f *= -1
	}
	return Transaction{
		ID:          t.ID,
		Description: t.Description,
		Type:        tpe,
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		// TODO convert date to relative date (ex. today, yesterday, ...)
		Date: t.Date.Format(time.DateOnly),
	}
}
