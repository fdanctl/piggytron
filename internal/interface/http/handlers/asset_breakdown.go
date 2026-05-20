package handlers

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type AssetBreakdownChartHandler struct{}

func NewAssetBreakdownChartHandler() *AssetBreakdownChartHandler {
	return &AssetBreakdownChartHandler{}
}

func (h *AssetBreakdownChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AssetBreakdownChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	chart := pieRadius()
	// cut unnecessary html code from echarts
	// only get what's inside <body></body>
	buf := bytes.NewBuffer(nil)
	chart.Render(buf)
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
	itemCntPie = 4
	seasons    = []string{"Spring", "Summer", "Autumn ", "Winter"}
)

func generatePieItems() []opts.PieData {
	items := make([]opts.PieData, 0)
	for i := 0; i < itemCntPie; i++ {
		items = append(items, opts.PieData{Name: seasons[i], Value: rand.Intn(100)})
	}
	return items
}

func pieRadius() *charts.Pie {
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

	pie.AddSeries("pie", generatePieItems()).
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
