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
	id transaction.ID,
	description string,
	ttype transaction.Ttype,
	amount int,
	date time.Time,
) Transaction {
	f := float64(amount) / 100
	if ttype == "expense" {
		f *= -1
	}
	return Transaction{
		Id:          id,
		Description: description,
		Type:        ttype,
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		// TODO convert date to relative date (ex. today, yesterday, ...)
		Date: date.Format(time.DateOnly),
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
