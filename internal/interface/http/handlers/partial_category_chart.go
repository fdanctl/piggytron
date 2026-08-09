package handlers

import (
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/views/charts"
)

type CategoryChartHandler struct {
	categoryQuery query.CategoryQueryService
}

func NewCategoryChartHandler(
	cq query.CategoryQueryService,
) *CategoryChartHandler {
	return &CategoryChartHandler{
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

	name := "Category"
	if len(mvalues) > 0 {
		name = mvalues[0].Name
	}

	theme := r.Header.Get("theme")
	bItems := charts.MakeCatBarItems(mvalues)
	chart := charts.CreateMonthlyBarChart(bItems, name, theme)
	charts.ConvertChartToTemplComponent(chart).Render(r.Context(), w)
}
