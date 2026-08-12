package charts

import (
	"math"
	"sort"
	"time"

	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// pair is used to sort a map by value
type pair struct {
	Key   string
	Value int
}

// GenerateAccountsHistLine computes, per account, the running total of
// the daily changes in hist and returns the resulting time series, plus the
// rounded min and max y-axis bounds for the chart.
func GenerateAccountsHistLine(
	hist []query.AccountDailyBalance,
) (map[string][]opts.LineData, int, int) {
	var min, max float64
	datas := make(map[string][]opts.LineData)

	for _, v := range hist {
		amount := float64(v.Balance) / float64(100)

		data := opts.LineData{Value: []any{v.Day, amount}}

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
	}

	rmin := 500 * ((int(math.Floor(min)) / 500) + 1) // previous multiple of 500
	rmax := 500 * ((int(math.Ceil(max)) / 500) + 1)  // next multiple of 500

	return datas, int(math.Min(0, float64(rmin))), rmax
}

// GenerateAccountsHistLineAndPieItems computes the running total of daily changes
// for each bank account, with all goals represented as a single aggregate. It returns
// the line chart data along with its keys sorted by balance, and rounded minimum
// and maximum y-axis bounds. It also returns pie chart data for the top 4 accounts
// by balance, with all remaining accounts grouped into an "Other" aggregate.
func GenerateAccountsHistLineAndPieItems(
	hist []query.AccountDailyBalance,
) (lineData map[string][]opts.LineData, sortedKeys []string, pieData []opts.PieData, ymin int, ymax int) {
	var min, max float64
	lineData = make(map[string][]opts.LineData)
	balances := make(map[string]int)

	date := hist[0].Day
	goalsK := "Goals"
	lineData[goalsK] = []opts.LineData{}
	var currGoalBalance int
	for _, v := range hist {
		// add Goals day balance
		if date != v.Day {
			value, _ := lineData[goalsK]
			amount := float64(currGoalBalance) / float64(100)
			data := opts.LineData{Value: []any{date, amount}}

			lineData[goalsK] = append(value, data)
			if amount < min {
				min = amount
			}
			if amount > max {
				max = amount
			}

			date = v.Day
			currGoalBalance = 0
		}
		if v.Type == "goal" {
			currGoalBalance += v.Balance
			continue
		}
		balances[v.Name] = v.Balance

		amount := float64(v.Balance) / float64(100)

		data := opts.LineData{Value: []any{v.Day, amount}}

		value, ok := lineData[v.Name]
		if !ok {
			lineData[v.Name] = []opts.LineData{data}
		} else {
			lineData[v.Name] = append(value, data)
		}

		if amount < min {
			min = amount
		}
		if amount > max {
			max = amount
		}
	}
	// add last day Goals balance
	balances[goalsK] = currGoalBalance
	amount := float64(currGoalBalance) / float64(100)
	data := opts.LineData{Value: []any{date, amount}}
	value, _ := lineData[goalsK]
	lineData[goalsK] = append(value, data)
	if amount < min {
		min = amount
	}
	if amount > max {
		max = amount
	}

	ymin = 500 * ((int(math.Floor(min)) / 500) + 1) // previous multiple of 500
	ymin = int(math.Min(0, float64(ymin)))
	ymax = 500 * ((int(math.Ceil(max)) / 500) + 1) // next multiple of 500

	pairs := make([]pair, 0, len(balances))
	for k, v := range balances {
		pairs = append(pairs, pair{k, v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Value > pairs[j].Value
	})

	availableSlots := 4
	var other int
	for i, p := range pairs {
		if i < availableSlots {
			sortedKeys = append(sortedKeys, p.Key)
			pieData = append(
				pieData,
				opts.PieData{Name: p.Key, Value: float64(p.Value) / float64(100)},
			)
		} else {
			other += p.Value
			// remove extra from line data
			delete(lineData, p.Key)
		}
	}

	pieData = append(pieData, opts.PieData{Name: "Other", Value: float64(other) / float64(100)})

	return
}

// LineTime builds a smooth multi-series line chart with a time x axis and
// its y axis bounded by [min, max].
func LineTime(
	m map[string][]opts.LineData,
	sortedKeys *[]string,
	min, max int,
	since time.Time,
	theme string,
) *charts.Line {
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
			"#50808E", "#C95D63", "#EEC170", "#A083A2", "#9B9B9B",
		}),
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

	if sortedKeys != nil {
		for _, k := range *sortedKeys {
			series, _ := m[k]
			line.AddSeries(k, series).
				SetSeriesOptions(
					charts.WithLineChartOpts(opts.LineChart{
						Smooth: opts.Bool(true),
						Symbol: "none",
					}),
				)
		}
	} else {
		for k, v := range m {
			line.AddSeries(k, v).
				SetSeriesOptions(
					charts.WithLineChartOpts(opts.LineChart{
						Smooth: opts.Bool(true),
						Symbol: "none",
					}),
				)
		}
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
				charts.WithLineChartOpts(opts.LineChart{
					Smooth: opts.Bool(true),
					Symbol: "none",
				}),
				charts.WithAreaStyleOpts(
					opts.AreaStyle{
						Opacity: opts.Float(0.5),
					}),
			)
	}
	return line
}
