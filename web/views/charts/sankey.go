package charts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"golang.org/x/text/currency"
)

// MakeSankey builds a sankey diagram from the given nodes and links, with
// optional entry animation.
func MakeSankey(
	sankeyNodes []opts.SankeyNode,
	sankeyLinks []opts.SankeyLink,
	animation bool,
	theme string,
) *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithAnimation(animation),
		charts.WithInitializationOpts(
			opts.Initialization{
				Width:           "100%",
				Height:          "100%",
				Theme:           theme,
				BackgroundColor: "transparent",
			},
		),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.8)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts(sankeyTooltipFormatter(currency.EUR)),
		}),
	)

	sankey.AddSeries(
		"sankey",
		sankeyNodes,
		sankeyLinks,
		charts.WithLineStyleOpts(opts.LineStyle{
			Color:     "target",
			Curveness: 0.5,
		}),
		charts.WithLabelOpts(opts.Label{
			Show: opts.Bool(true),
		}),
	)
	return sankey
}
