package views

import (
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/query"
)

// BudgetPageView is the view-model for the budget page: per-category rows
// grouped by needs/wants/savings, the summary totals and the percentage of
// the budget assigned to each group.
type BudgetPageView struct {
	Month               budget.Month
	TotalBudgeted       int
	ReadyToAssign       int
	Income              int
	AvailableToSpend    int
	Overspent           int
	NeedsRows           []BudgetRowView
	WantsRows           []BudgetRowView
	SavingsRows         []BudgetRowView
	CategoriesCarryover int
	UnassignCarryover   int

	NeedsAvailable int
	NeedsBudget    int
	NeedsPct       float64

	WantsAvailable int
	WantsBudget    int
	WantsPct       float64

	SavingsAvailable int
	SavingsBudget    int
	SavingsPct       float64
}

// BudgetRowView is the view-model for one category's budget row: budgeted
// amount, carryover from the previous month and available balance.
type BudgetRowView struct {
	CategoryID string
	Month      budget.Month
	Name       string
	Carryover  int
	Budgeted   int
	Available  int
}

// NewBudgetPageView aggregates the per-category budget values into rows and
// computes the month's totals: ready-to-assign, income, budgeted, spent, available,
// overspent and the carryover amounts.
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

	var needsAvailable int
	var wantsAvailable int
	var savingsAvailable int

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

		if v.ArchivedAt == nil || month.Time().Compare(*v.ArchivedAt) <= 0 {
			switch v.Type {
			case "needs":
				needs = append(needs, row)
				needsBudget += v.Budgeted
				needsAvailable += available
			case "wants":
				wants = append(wants, row)
				wantsBudget += v.Budgeted
				wantsAvailable += available
			case "savings":
				savings = append(savings, row)
				savingsBudget += v.Budgeted
				savingsAvailable += available
			}
		}
	}

	var needsPct float64
	if totalBudgeted > 0 {
		needsPct = (float64(needsBudget) / float64(totalBudgeted)) * 100
	}

	var wantsPct float64
	if totalBudgeted > 0 {
		wantsPct = (float64(wantsBudget) / float64(totalBudgeted)) * 100
	}

	var savingsPct float64
	if totalBudgeted > 0 {
		savingsPct = (float64(savingsBudget) / float64(totalBudgeted)) * 100
	}

	// last month budgeted but not spent
	categoriesCarryover := totalBudgetedPrev - totalSpentPrev

	return BudgetPageView{
		Month:               month,
		TotalBudgeted:       totalBudgeted,
		ReadyToAssign:       balance - totalAvailable,
		Income:              income,
		AvailableToSpend:    totalAvailable,
		Overspent:           overspent * -1,
		NeedsRows:           needs,
		WantsRows:           wants,
		SavingsRows:         savings,
		CategoriesCarryover: categoriesCarryover,
		UnassignCarryover:   max(balance-income-categoriesCarryover, 0),

		NeedsBudget:    needsBudget,
		NeedsAvailable: needsAvailable,
		NeedsPct:       needsPct,

		WantsBudget:    wantsBudget,
		WantsAvailable: wantsAvailable,
		WantsPct:       wantsPct,

		SavingsBudget:    savingsBudget,
		SavingsAvailable: savingsAvailable,
		SavingsPct:       savingsPct,
	}
}
