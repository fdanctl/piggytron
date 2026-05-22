package handlers

import (
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/application/charts"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
)

type AccountChartHandler struct {
	chartsService *charts.Service
	accountQuery  query.AccountQueryService
}

func NewAccountChartHandler(
	cs *charts.Service,
	aq query.AccountQueryService,
) *AccountChartHandler {
	return &AccountChartHandler{
		chartsService: cs,
		accountQuery:  aq,
	}
}

func (h *AccountChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if id == "" {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AccountChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")

	changeHist, err := h.accountQuery.GetAccountDailyChange(r.Context(), id)
	if err != nil {
		logger.Error("error finding accounts history", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	d := q.Get("start")
	m := q.Get("max")
	logger.Debug(m)

	int, err := strconv.Atoi(m)
	if err != nil {
		logger.Error("error parsing date", "error", err)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse(time.DateOnly, d)
	if err != nil {
		logger.Error("error parsing date", "error", err)
		http.Error(w, "invalid date", http.StatusBadRequest)
		return
	}

	histMap, _, max := h.chartsService.GenerateYearAccountsHistLine(changeHist)
	line := h.chartsService.LineTimeAccount(
		histMap,
		0,
		math.Max(float64(int)/100, float64(max)),
		startDate,
	)
	h.chartsService.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
