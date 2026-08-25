package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// CategoryChartHandler renders the monthly values bar chart for one
// category across the current year.
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

// Get renders the category's monthly bar chart for the current year.
func (h *CategoryChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")
	logger.Debug(id)

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "6m" // default
	}

	var startDate time.Time
	now := time.Now()
	switch period {
	case "1y":
		startDate = time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, now.Location())
	case "ytd":
		startDate = time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
	case "6m":
		startDate = time.Date(now.Year(), now.Month()-6, 1, 0, 0, 0, 0, now.Location())
	}

	mvalues, err := h.categoryQuery.GetMonthlyValueSince(r.Context(), id, startDate)
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
	bItems, xAxis := charts.MakeCatBarItems(mvalues)
	chart := charts.CreateMonthlyBarChart(bItems, name, xAxis, theme)
	components.SelectPill(
		"",
		strings.ToUpper(period),
		period,
		components.SelectPillLeft,
		[]components.SelectOption{
			{Label: "6M", Value: "6m"},
			{Label: "YTD", Value: "ytd"},
			{Label: "1Y", Value: "1y"},
		},
		templ.Attributes{
			"name":      "period",
			"hx-get":    fmt.Sprint("/partials/charts/cat-hist/", id),
			"hx-target": "closest .chart-container",
		},
	).Render(r.Context(), w)
	charts.ConvertChartToTemplComponent(chart).Render(r.Context(), w)
}
