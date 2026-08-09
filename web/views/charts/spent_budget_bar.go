package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// MakeBudgetSpentBarItems splits the category budget/spent data into the
// budget and spent series (excluding income categories) and the x-axis
// category names.
func MakeBudgetSpentBarItems(
	data []query.CategoryBudgetValue,
) (budget, spent []opts.BarData, xAxis []string) {
	for _, v := range data {
		if v.Type == "income" {
			continue
		}
		budget = append(budget, opts.BarData{Value: float64(v.Budgeted) / 100})
		spent = append(spent, opts.BarData{Value: float64(v.Value) / 100})
		xAxis = append(xAxis, v.Name)
	}

	return
}

// CreateCategoryBudgetSpentBarChart builds a grouped bar chart comparing the
// budgeted and spent amounts for the given categories.
func CreateCategoryBudgetSpentBarChart(
	budget, spent []opts.BarData,
	categories []string,
	theme string,
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
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(barTooltipFormatter(currency.EUR)),
		}),
		charts.WithXAxisOpts(opts.XAxis{
			AxisLabel: &opts.AxisLabel{
				Show: opts.Bool(false),
			},
		}),
	)

	bar.SetXAxis(categories).
		AddSeries("Budget", budget).
		AddSeries("Spent", spent)
	return bar
}
