package views

import (
	"fmt"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func formatMoney(amount float64, cur currency.Unit, lang language.Tag) string {
	p := message.NewPrinter(lang)
	symbol := currency.Symbol(cur)

	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return fmt.Sprintf(p.Sprintf("%s%s%.2f", sign, symbol, amount))
}
