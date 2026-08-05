package views

import (
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/query"
)

type BankPage struct {
	NetWorth       int
	NetWorthIncPct float64

	Savings       int
	SavingsIncPct float64

	Available       int
	AvailableIncPct float64

	Goals       int
	GoalsIncPct float64

	Banks              []Bank
	RecentTransactions []Transaction
}

type Bank struct {
	ID        string
	Name      string
	Type      string
	IsSavings bool
	Amount    int
}

func NewBank(id string, name string, btype string, isSaving bool, amount int) Bank {
	return Bank{
		ID:        id,
		Name:      name,
		Type:      btype,
		IsSavings: isSaving,
		Amount:    amount,
	}
}

func NewBankPage(
	a []query.AccountWithSumAndMonthChange,
	t []query.LedgerEntryDTO,
) BankPage {
	var savings int
	var savingsDelta int
	var available int
	var availableDelta int
	var goals int
	var goalsDelta int

	var banks []Bank
	for _, v := range a {
		var isSaving bool
		if v.Type == "bank" {
			if *v.IsSaving {
				savings += v.Sum
				savingsDelta += v.MoneyIn
				savingsDelta -= v.MoneyOut
				isSaving = true
			} else {
				available += v.Sum
				availableDelta += v.MoneyIn
				availableDelta -= v.MoneyOut
			}

			banks = append(
				banks,
				NewBank(v.ID, v.Name, v.Type, isSaving, v.Sum),
			)
		}

		if v.Type == string(account.GoalType) {
			goals += v.Sum
			goalsDelta += v.MoneyIn
			goalsDelta -= v.MoneyOut
		}
	}

	var tviews []Transaction
	for _, v := range t {
		tviews = append(tviews, NewTransaction(v))
	}

	nw := savings + available + goals
	nwDelta := savingsDelta + availableDelta + goalsDelta

	// increase percentages = delta / last month balance
	prevNW := nw - nwDelta
	var nwPct float64
	if prevNW >= 0 {
		nwPct = float64(nwDelta) / float64(prevNW) * 100
	}

	prevSavings := savings - savingsDelta
	var savingsPct float64
	if prevSavings >= 0 {
		savingsPct = float64(savingsDelta) / float64(prevSavings) * 100
	}

	prevAvailable := available - availableDelta
	var availablePct float64
	if prevAvailable >= 0 {
		availablePct = float64(availableDelta) / float64(prevAvailable) * 100
	}

	prevGoals := goals - goalsDelta
	var goalsPct float64
	if prevGoals >= 0 {
		goalsPct = float64(goalsDelta) / float64(prevGoals) * 100
	}

	return BankPage{
		NetWorth:           nw,
		NetWorthIncPct:     nwPct,
		Savings:            savings,
		SavingsIncPct:      savingsPct,
		Available:          available,
		AvailableIncPct:    availablePct,
		Goals:              goals,
		GoalsIncPct:        goalsPct,
		Banks:              banks,
		RecentTransactions: tviews,
	}
}
