package appcharts

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (s *Service) MakeBudgetSankeyNodeLink(
	name, catType string,
	value int,
) (opts.SankeyNode, opts.SankeyLink) {
	var src string
	var dst string
	var color string
	switch catType {
	case "income":
		src = name
		dst = "Budget"
		color = "#3a9a9a"
	case "needs":
		src = "Budget"
		dst = name
		color = "#95bf98"
	case "wants":
		src = "Budget"
		dst = name
		color = "#d9725b"
	case "savings":
		src = "Budget"
		dst = name
		color = "#bea9ba"
	}
	return opts.SankeyNode{
			Name: name,
			ItemStyle: &opts.ItemStyle{
				Color: color,
			},
		},
		opts.SankeyLink{
			Source: src,
			Target: dst,
			Value:  float32(value) / float32(100),
		}
}

func (s *Service) MakeSankey(
	sankeyNodes []opts.SankeyNode,
	sankeyLinks []opts.SankeyLink,
	animation bool,
) *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
		charts.WithAnimation(animation),
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			// Formatter:       opts.FuncOpts("myTooltipFormatter"),
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
			// TODO different colors for dark and light themes
			Color: "#fff",
		}),
	)
	return sankey
}
