package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BudgetChartHandler struct {
	chartsService *appcharts.Service
	categoryQuery query.CategoryQueryService
}

func NewBudgetChartHandler(
	cs *appcharts.Service,
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
		httperror.SendError(w, r, err)
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
		err := fmt.Errorf("error geting categories budget-spent: %w", err)
		httperror.SendError(w, r, err)
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

	component := components.NoData()
	if len(links) > 0 {
		sankey := h.chartsService.MakeSankey(nodes, links, true)
		component = h.chartsService.ConvertChartToTemplComponent(sankey)
	}

	component.Render(r.Context(), w)
}
