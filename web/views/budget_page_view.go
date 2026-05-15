package views

import (
	"github.com/fdanctl/piggytron/internal/query"
)

type BudgetPageView struct {
	TotalBudgeted int
	LeftToBudget  int
	Income        int
	LeftToSpend   int
	Overspent     int
	NeedsRows     []BudgetRowView
	WantsRows     []BudgetRowView
	SavingsRows   []BudgetRowView

	NeedsSpent  int
	NeedsBudget int
	NeedsPct    float64

	WantsSpent  int
	WantsBudget int
	WantsPct    float64

	SavingsSpent  int
	SavingsBudget int
	SavingsPct    float64
}

type BudgetRowView struct {
	CID      string
	BID      string
	Name     string
	Budgeted int
	Left     int
}

func NewBudgetPageView(
	income int,
	catBudgetSpent []query.ExpenseCategoryBudgetSpent,
) BudgetPageView {
	var totalBudgeted int
	var totalSpent int
	var overspent int

	var needsBudget int
	var wantsBudget int
	var savingsBudget int

	var needsSpent int
	var wantsSpent int
	var savingsSpent int

	var needs, wants, savings []BudgetRowView

	for _, v := range catBudgetSpent {
		totalBudgeted += v.Budgeted
		totalSpent += v.Spent
		left := v.Budgeted - v.Spent
		if left < 0 {
			overspent += left
		}

		row := BudgetRowView{
			CID:      v.CID,
			BID:      v.BID,
			Name:     v.Name,
			Budgeted: v.Budgeted,
			Left:     left,
		}

		switch v.Type {
		case "needs":
			needs = append(needs, row)
			needsBudget += v.Budgeted
			needsSpent += v.Spent
		case "wants":
			wants = append(wants, row)
			wantsBudget += v.Budgeted
			wantsSpent += v.Spent
		case "savings":
			savings = append(savings, row)
			savingsBudget += v.Budgeted
			savingsSpent += v.Spent
		}
	}

	return BudgetPageView{
		TotalBudgeted: totalBudgeted,
		LeftToBudget:  income - totalBudgeted,
		Income:        income,
		LeftToSpend:   income - totalSpent,
		Overspent:     overspent * -1,
		NeedsRows:     needs,
		WantsRows:     wants,
		SavingsRows:   savings,

		NeedsBudget: needsBudget,
		NeedsSpent:  needsSpent * -1,
		NeedsPct:    (float64(needsBudget) / float64(totalBudgeted)) * 100,

		WantsBudget: wantsBudget,
		WantsSpent:  wantsSpent * -1,
		WantsPct:    (float64(wantsBudget) / float64(totalBudgeted)) * 100,

		SavingsBudget: savingsBudget,
		SavingsSpent:  savingsSpent * -1,
		SavingsPct:    (float64(savingsBudget) / float64(totalBudgeted)) * 100,
	}
}
