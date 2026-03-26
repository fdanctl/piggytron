package views

import (
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
