package handlers

import (
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/application/charts"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
)

type CatHistChartHandler struct {
	chartsService *charts.Service
	categoryQuery query.CategoryQueryService
}

func NewCatHistChartHandler(
	cs *charts.Service,
	cq query.CategoryQueryService,
) *CatHistChartHandler {
	return &CatHistChartHandler{
		chartsService: cs,
		categoryQuery: cq,
	}
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
