package handlers

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

type DashboardBudgetCharts struct {
	chartsService *appcharts.Service
	categoryQuery query.CategoryQueryService
}

func NewDashboardBudgetCharts(
	cs *appcharts.Service,
	cq query.CategoryQueryService,
) *DashboardBudgetCharts {
	return &DashboardBudgetCharts{
		chartsService: cs,
		categoryQuery: cq,
	}
}

func (h *DashboardBudgetCharts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DashboardBudgetCharts) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	bm := budget.NewMonth(time.Now())
	categoryBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	pieItems := h.chartsService.MakeAssetsPieSpentItems(categoryBudgetSpent.Data)

	pie := components.NoData()
	if len(pieItems) > 0 {
		c := h.chartsService.PieRadius(pieItems)
		pie = h.chartsService.ConvertChartToTemplComponent(c)
	}

	budgetItems, spentItems, categories := h.chartsService.MakeBudgetSpentBarItems(
		categoryBudgetSpent.Data,
	)
	bar := h.chartsService.ConvertChartToTemplComponent(
		h.chartsService.CreateCategoryBudgetSpentBarChart(budgetItems, spentItems, categories),
	)

	templ.Join(
		bar,
		layouts.OOBWraper(
			"spent-by-type",
			"innerHTML",
			nil,
			pie,
		),
	).Render(r.Context(), w)
}
