package charts

import (
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeBudgetSpentBarItems splits the category budget/spent data into the
// budget and spent series (excluding income and archived categories) and the x-axis
// category names.
func MakeBudgetSpentBarItems(
	data []query.CategoryBudgetValue,
	month budget.Month,
) (budget, spent []opts.BarData, xAxis []string) {
	for _, v := range data {
		if v.Type == "income" || (v.ArchivedAt != nil && month.Time().After(*v.ArchivedAt)) {
			continue
		}
		budget = append(budget, opts.BarData{Value: float64(v.Budgeted) / 100})
		spent = append(spent, opts.BarData{Value: float64(v.Value) / 100})
		xAxis = append(xAxis, v.Name)
	}

	return
}
