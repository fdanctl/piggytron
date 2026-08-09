package charts

import (
	"slices"

	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// MakeAccountsPieItems converts account balances into pie chart data,
// sorted from largest to smallest balance.
func MakeAccountsPieItems(acc []query.AccountWithSum) []opts.PieData {
	slices.SortFunc(acc, func(a, b query.AccountWithSum) int {
		return b.Sum - a.Sum
	})
	var data []opts.PieData
	for _, v := range acc {
		data = append(data, opts.PieData{Name: v.Name, Value: float64(v.Sum) / float64(100)})
	}
	return data
}

// PieRadius builds a donut chart (40%-75% radius) from the given pie data.
func PieRadius(items []opts.PieData, name, theme string) *charts.Pie {
	pie := charts.NewPie()

	legentTextColor := "#afb3b3"
	if theme == "light" {
		legentTextColor = "#666966"
	}

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		// charts.WithLegendOpts(opts.Legend{
		// 	Show:   opts.Bool(false),
		// 	Top:    "center",
		// 	Right:  "0",
		// 	Orient: "vertical",
		// }),
		charts.WithColorsOpts(opts.Colors{
			"#95bf98", "#d9725b", "#b185a7", "#297373",
		}),
		charts.WithLegendOpts(opts.Legend{
			TextStyle: &opts.TextStyle{
				Color: legentTextColor,
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(pieTooltipFormatter(currency.EUR)),
		}),
	)

	pie.AddSeries(name, items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show: opts.Bool(false),
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}
