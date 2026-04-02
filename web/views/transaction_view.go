package views

import (
	"github.com/fdanctl/piggytron/internal/query"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

type Transaction struct {
	ID          string
	Description string
	Type        string
	Category    string
	Accounts    []string
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

	cat := "Transfer"
	if t.ExpenseCategory != nil {
		cat = *t.ExpenseCategory
	} else if t.IncomeCategory != nil {
		cat = *t.IncomeCategory
	}

	accs := make([]string, 0, 2)
	if t.FromAccount != nil {
		accs = append(accs, *t.FromAccount)
	}
	if t.ToAccount != nil {
		accs = append(accs, *t.ToAccount)
	}

	return Transaction{
		ID:          t.ID,
		Description: t.Description,
		Type:        t.Type,
		Category:    cat,
		Accounts:    accs,
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		Date:        FormatDate(t.Date),
	}
}

func NewAccountTransaction(
	t query.TransactionDTO,
	a *query.AccountWithCategory,
) Transaction {
	f := float64(t.Amount) / 100
	tpe := "income"
	if t.FromAccount != nil && *t.FromAccount == a.Name {
		tpe = "expense"
		f *= -1
	}

	var cat string
	if t.ExpenseCategory != nil {
		cat = *t.ExpenseCategory
	} else {
		cat = *t.IncomeCategory
	}

	accs := make([]string, 0, 2)
	if t.FromAccount != nil {
		accs = append(accs, *t.FromAccount)
	}
	if t.ToAccount != nil {
		accs = append(accs, *t.ToAccount)
	}

	return Transaction{
		ID:          t.ID,
		Description: t.Description,
		Type:        tpe,
		Category:    cat,
		Accounts:    accs,
		Amount:      formatMoney(f, currency.EUR, language.AmericanEnglish),
		Date:        FormatDate(t.Date),
	}
}
