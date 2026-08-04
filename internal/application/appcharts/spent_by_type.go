package appcharts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (s *Service) MakeAssetsPieSpentItems(cats []query.CategoryBudgetValue) []opts.PieData {
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
