package views

import (
	"time"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Transaction struct {
	Id          transaction.ID
	Description string
	Type        transaction.Ttype
	Amount      string
	Date        string
}

func NewTransaction(
	t *transaction.Transaction,
) Transaction {
	f := float64(t.Amount()) / 100
	if t.Ttype() == "expense" {
		f *= -1
	}
	return Transaction{
		Id:          t.ID(),
		Description: t.Description(),
		Type:        t.Ttype(),
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		// TODO convert date to relative date (ex. today, yesterday, ...)
		Date: t.Date().Format(time.DateOnly),
	}
}
