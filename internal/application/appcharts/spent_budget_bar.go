package appcharts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (s *Service) MakeBudgetSpentBarItems(
	data []query.CategoryBudgetValue,
) (budget, spent []opts.BarData, xAxis []string) {
	for _, v := range data {
		if v.Type == "income" {
			continue
		}
		budget = append(budget, opts.BarData{Value: v.Budgeted})
		spent = append(spent, opts.BarData{Value: v.Value})
		xAxis = append(xAxis, v.Name)
	}

	return
}

func (s *Service) CreateCategoryBudgetSpentBarChart(
	budget, spent []opts.BarData,
	categories []string,
) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithColorsOpts(opts.Colors{
			"#5eefef", "#297373",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			// Formatter:       opts.FuncOpts("myTooltipFormatter"),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show: opts.Bool(false),
			},
		}),
	)

	bar.Assets.ClearPresetJSAssets()
	bar.SetXAxis(categories).
		AddSeries("Budget", budget).
		AddSeries("Spent", spent)
	return bar
}
