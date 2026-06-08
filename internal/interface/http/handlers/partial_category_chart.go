package handlers

import (
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
)

type CategoryChartHandler struct {
	chartsService *appcharts.Service
	categoryQuery query.CategoryQueryService
}

func NewCategoryChartHandler(
	cs *appcharts.Service,
	cq query.CategoryQueryService,
) *CategoryChartHandler {
	return &CategoryChartHandler{
		chartsService: cs,
		categoryQuery: cq,
	}
}

func (h *CategoryChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoryChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")
	logger.Debug(id)

	mvalues, err := h.categoryQuery.GetYearMonthlyValue(r.Context(), time.Now().Year(), id)
	if err != nil {
		logger.Error("error finding values", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bItems := h.chartsService.MakeBarItems(mvalues)
	chart := h.chartsService.CreateMonthlyBarChart(bItems)
	h.chartsService.ConvertChartToTemplComponent(chart).Render(r.Context(), w)
}
