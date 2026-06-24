package views

import (
	"math"

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

	NeedsLeft   int
	NeedsBudget int
	NeedsPct    float64

	WantsLeft   int
	WantsBudget int
	WantsPct    float64

	SavingsLeft   int
	SavingsBudget int
	SavingsPct    float64
}

type BudgetRowView struct {
	CategoryID string
	BudgetID   string
	Name       string
	Budgeted   int
	Left       int
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
			CategoryID: v.CategoryID,
			BudgetID:   v.BudgetID,
			Name:       v.Name,
			Budgeted:   v.Budgeted,
			Left:       left,
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

	needsPct := (float64(needsBudget) / float64(totalBudgeted)) * 100
	wantsPct := (float64(wantsBudget) / float64(totalBudgeted)) * 100
	savingsPct := (float64(savingsBudget) / float64(totalBudgeted)) * 100
	if math.IsNaN(needsPct) {
		needsPct = 0
	}
	if math.IsNaN(wantsPct) {
		wantsPct = 0
	}
	if math.IsNaN(savingsPct) {
		savingsPct = 0
	}

	return BudgetPageView{
		TotalBudgeted: totalBudgeted,
		LeftToBudget:  income - totalBudgeted,
		Income:        income,
		LeftToSpend:   totalBudgeted - totalSpent,
		Overspent:     overspent * -1,
		NeedsRows:     needs,
		WantsRows:     wants,
		SavingsRows:   savings,

		NeedsBudget: needsBudget,
		NeedsLeft:   needsBudget - needsSpent,
		NeedsPct:    needsPct,

		WantsBudget: wantsBudget,
		WantsLeft:   wantsBudget - wantsSpent,
		WantsPct:    wantsPct,

		SavingsBudget: savingsBudget,
		SavingsLeft:   savingsBudget - savingsSpent,
		SavingsPct:    savingsPct,
	}
}
