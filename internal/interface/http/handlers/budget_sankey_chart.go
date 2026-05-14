package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BudgetChartHandler struct{}

func NewBudgetChartHandler() *BudgetChartHandler {
	return &BudgetChartHandler{}
}

func (h *BudgetChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BudgetChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	month := r.PathValue("month")
	logger := middleware.LoggerFromContext(r.Context())
	logger.Debug(month)

	line := sankeyBase()
	// cut unnecessary html code from echarts
	// only get what's inside <body></body>
	buf := bytes.NewBuffer(nil)
	line.Render(buf)
	html := buf.String()
	bodyStart := strings.Index(html, "<body>")
	bodyEnd := strings.Index(html, "</body>")
	fragment := html[bodyStart+len("<body>") : bodyEnd]

	idStart := strings.Index(fragment, "id=\"")
	id := fragment[idStart+len("id=\"") : idStart+len("id=\"")+12]

	fmt.Fprint(w, fragment)
	// ResizeObserver script
	fmt.Fprintf(w, `<script>
		// block resize when chart is animating
		let initialized_%s = false;
		const observer_%s = new ResizeObserver(() => {
			if (!initialized_%s) {
				initialized_%s = true;
				return;
			}
			goecharts_%s.resize()
		})
		observer_%s.observe(document.getElementById("%s"))
		</script>`, id, id, id, id, id, id, id,
	)
}

var (
	sankeyNode = []opts.SankeyNode{
		{Name: "Salary"},
		{Name: "Side hustle"},
		{
			Name: "Budget",
			ItemStyle: &opts.ItemStyle{
				Color: "#194e4e",
			},
		},

		{
			Name: "Rent",
			ItemStyle: &opts.ItemStyle{
				Color: "#95bf98", // needs
			},
		},
		{
			Name: "Utilities",
			ItemStyle: &opts.ItemStyle{
				Color: "#95bf98", // needs
			},
		},
		{
			Name: "Eating out",
			ItemStyle: &opts.ItemStyle{
				Color: "#d9725b", // wants
			},
		},
		{
			Name: "Investments",
			ItemStyle: &opts.ItemStyle{
				Color: "#bea9ba", // savings
			},
		},
	}

	sankeyLink = []opts.SankeyLink{
		{Source: "Salary", Target: "Budget", Value: 80},
		{Source: "Side hustle", Target: "Budget", Value: 20},
		{Source: "Budget", Target: "Rent", Value: 40},
		{Source: "Budget", Target: "Utilities", Value: 10},
		{Source: "Budget", Target: "Eating out", Value: 25},
		{Source: "Budget", Target: "Investments", Value: 25},
	}
)

func sankeyBase() *charts.Sankey {
	sankey := charts.NewSankey()
	sankey.SetGlobalOptions(
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
		sankeyNode,
		sankeyLink,
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
