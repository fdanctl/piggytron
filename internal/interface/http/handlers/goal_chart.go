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

type GoalChartHandler struct{}

func NewGoalChartHandler() *GoalChartHandler {
	return &GoalChartHandler{}
}

func (h *GoalChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GoalChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	line := lineTime()
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
	itemCntLine = 6
	fruits      = []string{"Apple", "Banana", "Peach ", "Lemon", "Pear", "Cherry"}
)

func generateLineItemsTwoAxis(points int, xFunc func(int) interface{}) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < points; i++ {
		items = append(items, opts.LineData{Value: []interface{}{xFunc(i), 100 + rand.Intn(20)}})
	}
	return items
}

func lineTime() *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{Width: "100%", Height: "100%"}),
		charts.WithLegendOpts(opts.Legend{
			Show: opts.Bool(false),
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Min: 0,
			Max: 200,
		}),
		charts.WithColorsOpts(opts.Colors{
			"#5eefef", "#4bc4c4",
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "time",
			Min:  time.Date(2025, time.January, 1, 0, 0, 0, 0, time.Local),
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Trigger:         "axis",
			BackgroundColor: "rgba(0, 0, 0, 0.7)",
			BorderColor:     "transparent",
			// Formatter:       opts.FuncOpts("myTooltipFormatter"),
		}),
	)

	line.AddSeries(
		"History",
		generateLineItemsTwoAxis(
			20,
			func(i int) interface{} { return time.Date(2025, time.January, i, 0, 0, 0, 0, time.Local) },
		),
	).SetSeriesOptions(
		charts.WithLineChartOpts(
			opts.LineChart{
				Smooth: opts.Bool(true),
			}),
		charts.WithAreaStyleOpts(
			opts.AreaStyle{
				Opacity: opts.Float(0.5),
			}),
	)
	return line
}
