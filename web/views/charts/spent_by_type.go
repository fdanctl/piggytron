package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeBudgetSpentByTypePieItems aggregates the budget and spent value per
// category type (needs, wants, savings) into two pie chart data, respectively.
func MakeBudgetSpentByTypePieItems(
	cats []query.CategoryBudgetValue,
) (spentData []opts.PieData, budgetedData []opts.PieData) {
	var needsSpent, wantsSpent, savingsSpent int
	var needsBudgeted, wantsBudgeted, savingsBudgeted int
	for _, v := range cats {
		switch v.Type {
		case "needs":
			needsSpent += v.Value
			needsBudgeted += v.Budgeted
		case "wants":
			wantsSpent += v.Value
			wantsBudgeted += v.Budgeted
		case "savings":
			savingsSpent += v.Value
			savingsBudgeted += v.Budgeted
		}
	}

	spentData = []opts.PieData{
		{
			Name:  "Needs",
			Value: float64(needsSpent) / 100,
		},
		{
			Name:  "Wants",
			Value: float64(wantsSpent) / 100,
		},
		{
			Name:  "Savings",
			Value: float64(savingsSpent) / 100,
		},
	}

	budgetedData = []opts.PieData{
		{
			Name:  "Needs",
			Value: float64(needsBudgeted) / 100,
		},
		{
			Name:  "Wants",
			Value: float64(wantsBudgeted) / 100,
		},
		{
			Name:  "Savings",
			Value: float64(savingsBudgeted) / 100,
		},
	}
	return
}
