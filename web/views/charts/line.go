package charts

import (
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

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
			series := m[k]
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
