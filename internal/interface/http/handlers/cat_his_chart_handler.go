package handlers

import (
	"bytes"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type CatHistChartHandler struct{}

func NewCatHistChartHandler() *CatHistChartHandler {
	return &CatHistChartHandler{}
}

func (h *CatHistChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CatHistChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	chart := createBarChart()
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
	fmt.Fprint(w, fmt.Sprintf(`<script>
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
		</script>`, id, id, id, id, id, id, id),
	)
}

func generateBarItems() []opts.BarData {
	items := make([]opts.BarData, 0)
	for range 12 {
		items = append(items, opts.BarData{Value: rand.Intn(300)})
	}
	return items
}

func createBarChart() *charts.Bar {
	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithColorsOpts(opts.Colors{
			"#5eefef", "#4bc4c4",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			Formatter:       opts.FuncOpts("myTooltipFormatter"),
		}),
	)

	abbv := []string{
		"Jan", "Feb", "Mar", "Apr", "May", "Jun",
		"Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
	}

	currMonth := time.Now().Month() - 1
	xAxis := make([]string, 12, 12)
	for i := 11; i >= 0; i-- {
		m := (((int(currMonth) - i) + 12) % 12)
		xAxis[11-i] = abbv[m]
	}

	bar.Assets.ClearPresetJSAssets()
	bar.SetXAxis(xAxis).
		AddSeries("Expense", generateBarItems())
	return bar
}
