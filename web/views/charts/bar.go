package charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

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

// CreateMonthlyBarChart builds a bar chart of the 12 months of the year
// from the given bar data.
func CreateMonthlyBarChart(items []opts.BarData, name, theme string) *charts.Bar {
	barColor := "#5eefef"
	if theme == "light" {
		barColor = "#4bc4c4"
	}
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithColorsOpts(opts.Colors{barColor}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(barTooltipFormatter(currency.EUR)),
		}),
	)

	abbv := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}

	bar.SetXAxis(abbv).
		AddSeries(name, items).
		SetSeriesOptions(charts.WithMarkLineNameTypeItemOpts(
			opts.MarkLineNameTypeItem{
				Name: "Average", Type: "average",
				LineStyle: &opts.LineStyle{
					Type:    "dashed",
					Opacity: opts.Float(0.75),
				},
			},
		),
			charts.WithMarkLineStyleOpts(opts.MarkLineStyle{
				Symbol: []string{"none", "none"},
				Label:  &opts.Label{Show: opts.Bool(false)},
			}),
		)
	return bar
}
