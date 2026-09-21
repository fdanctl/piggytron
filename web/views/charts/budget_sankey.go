package charts

import (
	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
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

func BudgetChart(
	unassignCarry, onHold int,
	theme string,
	categoriesBudget []query.CategoryBudget,
) templ.Component {
	nodes := []opts.SankeyNode{
		{
			Name: "Budget",
			ItemStyle: &opts.ItemStyle{
				Color: "#194e4e",
			},
		},
	}
	var links []opts.SankeyLink
	if unassignCarry > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "Unassigned Carryover",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Unassigned Carryover",
				Target: "Budget",
				Value:  float32(unassignCarry) / float32(100),
			},
		)
	}
	budget := unassignCarry
	var budgeted int
	for _, v := range categoriesBudget {
		if v.Type == "income" {
			budget += v.Value
		} else {
			budgeted += v.Value
		}
		if v.Value > 0 {
			node, link := MakeBudgetSankeyNodeLink(v.Name, v.Type, v.Value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	ltb := budget - budgeted - onHold
	if ltb > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "Unassigned",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Budget",
				Target: "Unassigned",
				Value:  float32(ltb) / float32(100),
			},
		)
	}

	if onHold > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "On hold",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Budget",
				Target: "On hold",
				Value:  float32(onHold) / float32(100),
			},
		)
	}

	component := components.NoData()
	if len(links) > 0 {
		sankey := MakeSankey(nodes, links, true, theme)
		component = ConvertChartToTemplComponent(sankey)
	}

	return component
}
