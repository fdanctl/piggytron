package charts

import (
	"math"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// GenerateYearAccountsHistLine computes, per account, the running total of
// the daily changes in hist and returns the resulting time series, plus the
// rounded min and max y-axis bounds for the chart.
func GenerateYearAccountsHistLine(
	hist []query.AccountDailyChange,
) (map[string][]opts.LineData, int, int) {
	var min, max float64
	datas := make(map[string][]opts.LineData)

	runningTotals := make(map[string]int)
	lastDate := make(map[string]time.Time)
	for _, v := range hist {
		total := runningTotals[v.Name] + v.Change
		runningTotals[v.Name] = total

		amount := float64(total) / float64(100)
		data := opts.LineData{Value: []any{v.Date, amount}}
		value, ok := datas[v.Name]
		if !ok {
			datas[v.Name] = []opts.LineData{data}
		} else {
			datas[v.Name] = append(value, data)
		}

		if amount < min {
			min = amount
		}
		if amount > max {
			max = amount
		}
		lastDate[v.Name] = v.Date
	}

	// add today
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	for k, v := range datas {
		if !today.Equal(lastDate[k]) {
			amount := float64(runningTotals[k]) / float64(100)
			datas[k] = append(v, opts.LineData{Value: []any{today, amount}})
		}
	}

	rmin := 500 * ((int(math.Floor(min)) / 500) + 1) // previous multiple of 500
	rmax := 500 * ((int(math.Ceil(max)) / 500) + 1)  // next multiple of 500

	return datas, int(math.Min(0, float64(rmin))), rmax
}

// LineTime builds a smooth multi-series line chart with a time x axis and
// its y axis bounded by [min, max].
func LineTime(m map[string][]opts.LineData, min, max int, theme string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: min,
			Max: max,
		}),
		charts.WithColorsOpts(opts.Colors{
			"#b185a7", "#95bf98", "#d9725b", "#297373",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "time",
			Min:  time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger:         "axis",
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(lineTooltipFormatter(currency.EUR)),
		}),
	)

	for k, v := range m {
		line.AddSeries(k, v).
			SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))
	}
	return line
}

// LineTimeAccount is like LineTime but accepts fractional y-axis bounds and
// draws each series as a filled area starting from the given date.
func LineTimeAccount(
	m map[string][]opts.LineData,
	min, max float64,
	since time.Time,
	theme string,
) *charts.Line {
	lineColor := "#5eefef"
	if theme == "light" {
		lineColor = "#4bc4c4"
	}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:           "100%",
				Height:          "100%",
				BackgroundColor: "transparent",
			},
		),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: min,
			Max: max,
		}),
		charts.WithColorsOpts(opts.Colors{lineColor}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "time",
			Min:  since,
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger:         "axis",
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(lineTooltipFormatter(currency.EUR)),
		}),
	)

	for k, v := range m {
		line.AddSeries(k, v).
			SetSeriesOptions(
				charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}),
				charts.WithAreaStyleOpts(
					opts.AreaStyle{
						Opacity: opts.Float(0.5),
					}),
			)
	}
	return line
}
