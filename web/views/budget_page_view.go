package views

import (
	"math"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/query"
)

type BudgetPageView struct {
	Month             budget.Month
	TotalBudgeted     int
	ReadyToAssign     int
	Income            int
	AvailableToSpend  int
	Overspent         int
	NeedsRows         []BudgetRowView
	WantsRows         []BudgetRowView
	SavingsRows       []BudgetRowView
	TotalCarryover    int
	UnassignCarryover int

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
	Month      budget.Month
	Name       string
	Carryover  int
	Budgeted   int
	Available  int
}

func NewBudgetPageView(
	month budget.Month,
	net int,
	balance int,
	catBudgetSpent []query.CategoryBudgetValue,
) BudgetPageView {
	var income int
	var totalBudgetedPrev int
	var totalSpentPrev int

	var totalBudgeted int
	var totalSpent int
	var totalAvailable int
	var overspent int

	var needsBudget int
	var wantsBudget int
	var savingsBudget int

	var needsSpent int
	var wantsSpent int
	var savingsSpent int

	var needs, wants, savings []BudgetRowView

	for _, v := range catBudgetSpent {
		if v.Type == "income" {
			income += v.Value
			continue
		}
		totalBudgetedPrev += v.PrevTotalBudget
		totalSpentPrev += v.PrevTotalSpent

		totalBudgeted += v.Budgeted
		totalSpent += v.Value
		carryover := v.PrevTotalBudget - v.PrevTotalSpent
		available := v.Budgeted - v.Value + carryover
		totalAvailable += available
		if available < 0 {
			overspent += available
		}

		row := BudgetRowView{
			CategoryID: v.CategoryID,
			Month:      budget.Month(v.Month),
			Name:       v.Name,
			Carryover:  carryover,
			Budgeted:   v.Budgeted,
			Available:  available,
		}

		switch v.Type {
		case "needs":
			needs = append(needs, row)
			needsBudget += v.Budgeted
			needsSpent += v.Value
		case "wants":
			wants = append(wants, row)
			wantsBudget += v.Budgeted
			wantsSpent += v.Value
		case "savings":
			savings = append(savings, row)
			savingsBudget += v.Budgeted
			savingsSpent += v.Value
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
		Month:             month,
		TotalBudgeted:     totalBudgeted,
		ReadyToAssign:     balance - totalAvailable,
		Income:            income,
		AvailableToSpend:  totalAvailable,
		Overspent:         overspent * -1,
		NeedsRows:         needs,
		WantsRows:         wants,
		SavingsRows:       savings,
		TotalCarryover:    totalBudgetedPrev - totalSpentPrev,
		UnassignCarryover: max((balance-net)-(totalBudgetedPrev-totalSpentPrev), 0),

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
