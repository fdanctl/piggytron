package views

import (
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/query"
)

type BankPage struct {
	NetWorth int
	// net worth last month comparition
	Savings int
	// savings last month comparition
	Available int
	// available last month comparition
	Goals int
	// last month comparition
	Banks              []Bank
	RecentTransactions []Transaction
}

type Bank struct {
	ID   string
	Name string
	Type string
	// Amount       string
}

func NewBankPage(a []query.AccountWithSum, t []query.LedgerEntryDTO) BankPage {
	var savings int
	var available int
	var goals int

	var banks []Bank
	for _, v := range a {
		if v.Type == "bank" {
			if *v.IsSaving {
				savings += v.Sum
			} else {
				available += v.Sum
			}

			banks = append(banks, Bank{ID: v.ID, Name: v.Name, Type: v.Type})
		}

		if v.Type == string(account.GoalType) {
			goals += v.Sum
		}
	}

	var tviews []Transaction
	for _, v := range t {
		tviews = append(tviews, NewTransaction(v))
	}

	return BankPage{
		NetWorth:           savings + available + goals,
		Savings:            savings,
		Available:          available,
		Goals:              goals,
		Banks:              banks,
		RecentTransactions: tviews,
	}
}
