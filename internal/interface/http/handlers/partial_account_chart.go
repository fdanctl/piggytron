package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
)

type AccountChartHandler struct {
	chartsService *appcharts.Service
	accountQuery  query.AccountQueryService
}

func NewAccountChartHandler(
	cs *appcharts.Service,
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
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	q := r.URL.Query()
	d := q.Get("start")
	m := q.Get("max")
	logger.Debug(m)

	qmax, err := strconv.Atoi(m)
	if err != nil {
		qmax = 0
	}

	startDate, err := time.Parse(time.DateOnly, d)
	if err != nil {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	}

	histMap, _, max := h.chartsService.GenerateYearAccountsHistLine(changeHist)
	line := h.chartsService.LineTimeAccount(
		histMap,
		0,
		math.Max(float64(qmax)/100, float64(max)),
		startDate,
	)
	h.chartsService.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
