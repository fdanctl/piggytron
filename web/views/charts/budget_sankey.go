package charts

import (
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeBudgetSankeyNodeLink maps a category of the given type to the node and
// link it contributes to the budget sankey diagram. Income flows into the
// budget, while needs, wants and savings flow out of it.
func MakeBudgetSankeyNodeLink(
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
