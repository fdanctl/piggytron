package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeSpentByTypePieItems aggregates the spent value per budget category
// type (needs, wants, savings) into pie chart data.
func MakeSpentByTypePieItems(cats []query.CategoryBudgetValue) []opts.PieData {
	var needs, wants, savings int
	for _, v := range cats {
		switch v.Type {
		case "needs":
			needs += v.Value
		case "wants":
			wants += v.Value
		case "savings":
			savings += v.Value
		}
	}

	data := []opts.PieData{
		{
			Name:  "Needs",
			Value: needs,
		},
		{
			Name:  "Wants",
			Value: wants,
		},
		{
			Name:  "Savings",
			Value: savings,
		},
	}
	return data
}
