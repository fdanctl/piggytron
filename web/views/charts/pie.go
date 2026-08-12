package charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// PieRadius builds a donut chart (40%-75% radius) from the given pie data.
func PieRadius(items []opts.PieData, name, theme string) *charts.Pie {
	pie := charts.NewPie()

	legentTextColor := "#afb3b3"
	if theme == "light" {
		legentTextColor = "#666966"
	}

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithColorsOpts(opts.Colors{
			"#50808E", "#C95D63", "#EEC170", "#A083A2", "#9B9B9B",
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

// PieInPie builds a two donut stacked chart.
func PieInPie(
	outerItems []opts.PieData, outerName string,
	innerItems []opts.PieData, innerName string,
	theme string,
) *charts.Pie {
	pie := charts.NewPie()

	legentTextColor := "#afb3b3"
	if theme == "light" {
		legentTextColor = "#666966"
	}

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithColorsOpts(opts.Colors{
			"#95bf98", "#d9725b", "#b185a7",
		}),
		charts.WithLegendOpts(opts.Legend{
			TextStyle: &opts.TextStyle{
				Color: legentTextColor,
			},
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(pieTooltipFormatterWithSeries(currency.EUR)),
		}),
	)

	pie.AddSeries(outerName, outerItems,
		charts.WithLabelOpts(opts.Label{
			Show: opts.Bool(false),
		}),
		charts.WithPieChartOpts(opts.PieChart{
			Radius: []string{"50%", "75%"},
		}),
	)

	pie.AddSeries(innerName, innerItems,
		charts.WithLabelOpts(opts.Label{
			Show: opts.Bool(false),
		}),
		charts.WithPieChartOpts(opts.PieChart{
			Radius: []string{"23%", "48%"},
		}),
	)
	return pie
}
