package handlers

import (
	"net/http"
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
	chartComponent := h.chartsService.ConvertChartToTemplComponent(sankey)
	chartComponent.Render(r.Context(), w)
}
