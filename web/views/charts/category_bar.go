package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeCatBarItems maps the per-month category values to a fixed 12-month
// slice of bar data (one entry per month, zero-filled for missing months).
func MakeCatBarItems(mvalues []query.CategoryMonthlyValue) []opts.BarData {
	byMonth := make(map[int]int, len(mvalues))
	for _, v := range mvalues {
		byMonth[v.Month] = v.Value
	}
	items := make([]opts.BarData, 0, 12)
	for m := 1; m <= 12; m++ {
		items = append(items, opts.BarData{Value: float64(byMonth[m]) / 100})
	}
	return items
}

// CreateMonthlyBarChart builds a bar chart of the 12 months of the year
// from the given bar data.
func CreateMonthlyBarChart(items []opts.BarData) *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithColorsOpts(opts.Colors{
			"#5eefef", "#4bc4c4",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			// Formatter:       opts.FuncOpts("myTooltipFormatter"),
		}),
	)

	abbv := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}

	bar.Assets.ClearPresetJSAssets()
	bar.SetXAxis(abbv).
		AddSeries("Value", items)
	return bar
}
