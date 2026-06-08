package appcharts

import (
	"slices"

	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func (s *Service) MakeAssetsPieItems(acc []query.AccountWithSum, count int) []opts.PieData {
	// var othersValue int
	slices.SortFunc(acc, func(a, b query.AccountWithSum) int {
		return b.Sum - a.Sum
	})
	var data []opts.PieData
	for _, v := range acc {
		data = append(data, opts.PieData{Name: v.Name, Value: float64(v.Sum) / float64(100)})
	}
	return data
}

func (s *Service) PieRadius(items []opts.PieData) *charts.Pie {
	pie := charts.NewPie()
	const formatterJS = `
		function  myTooltipFormatter(p) {
  var color = p.color || '#666';
  return '<span style="color:' + color + ';font-size:14px;font-weight:bold;">' +
         p.seriesName + '<br/>' +
         p.name + ': ' + p.value + 
         (p.percent ? ' (' + p.percent + '%)' : '') +
         '</span>';
}`

	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "300%", Height: "100%"}),
		// charts.WithLegendOpts(opts.Legend{
		// 	Show:   opts.Bool(false),
		// 	Top:    "center",
		// 	Right:  "0",
		// 	Orient: "vertical",
		// }),
		charts.WithColorsOpts(opts.Colors{
			"#b185a7", "#95bf98", "#d9725b", "#297373",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			// Formatter:       opts.FuncOpts("myTooltipFormatter"),
		}),
	)

	pie.AddSeries("pie", items).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{
				Show:      opts.Bool(false),
				Formatter: "{b}: {c}",
			}),
			charts.WithPieChartOpts(opts.PieChart{
				Radius: []string{"40%", "75%"},
			}),
		)
	return pie
}
