package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (s *Service) MakeBarItems(mvalues []query.CategoryMonthlyValue) []opts.BarData {
	if len(mvalues) == 0 {
		return []opts.BarData{}
	}

	items := make([]opts.BarData, 0, 12)
outerLoop:
	for m := 1; m <= 12; m++ {
		for i, v := range mvalues {
			if m == v.Month {
				items = append(items, opts.BarData{Value: float64(v.Value) / float64(100)})
				mvalues = append(mvalues[:i], mvalues[i+1:]...)
				continue outerLoop
			}
		}
		items = append(items, opts.BarData{Value: 0})
	}
	return items
}

func (s *Service) CreateMonthlyBarChart(items []opts.BarData) *charts.Bar {
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
