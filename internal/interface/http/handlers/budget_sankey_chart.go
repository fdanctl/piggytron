package handlers

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/fdanctl/piggytron/internal/application/charts"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BudgetChartHandler struct {
	chartsService *charts.Service
	categoryQuery query.CategoryQueryService
}

func NewBudgetChartHandler(
	cs *charts.Service,
	cq query.CategoryQueryService,
) *BudgetChartHandler {
	return &BudgetChartHandler{
		chartsService: cs,
		categoryQuery: cq,
	}
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
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	now := time.Now()

	minD := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	maxD := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	categoryBudget, err := h.categoryQuery.GetCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		minD,
		maxD,
	)
	if err != nil {
		logger.Error("error geting category budget-spent", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	nodes := []opts.SankeyNode{
		{
			Name: "Budget",
			ItemStyle: &opts.ItemStyle{
				Color: "#194e4e",
			},
		},
	}
	var links []opts.SankeyLink

	for _, v := range categoryBudget {
		if v.Value > 0 {
			node, link := h.chartsService.MakeBudgetSankeyNodeLink(v.Name, v.Type, v.Value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	sankey := h.chartsService.MakeSankey(nodes, links, true)

	// TODO MAKE THIS A REUSABLE FUNC
	// cut unnecessary html code from echarts
	// only get what's inside <body></body>
	buf := bytes.NewBuffer(nil)
	sankey.Render(buf)
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
