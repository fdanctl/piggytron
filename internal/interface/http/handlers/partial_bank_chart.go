package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type BankChartHandler struct {
	chartsService *appcharts.Service
	accountQuery  query.AccountQueryService
}

func NewBankChartHandler(
	cs *appcharts.Service,
	aq query.AccountQueryService,
) *BankChartHandler {
	return &BankChartHandler{
		chartsService: cs,
		accountQuery:  aq,
	}
}

func (h *BankChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *BankChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")

	changeHist, err := h.accountQuery.GetAccountDailyChange(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	q := r.URL.Query()
	month := q.Get("month")

	var startDate time.Time
	if month == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	} else {
		y, m, err := parseMonth(month)
		if err != nil {
			logger.Error("unexpected error", "error", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		logger.Debug("parseMonth", "year", y, "month", m)
		startDate = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
	}

	histMap, _, max := h.chartsService.GenerateYearAccountsHistLine(changeHist)
	line := h.chartsService.LineTimeAccount(
		histMap,
		0,
		float64(max),
		startDate,
	)
	chartComponent := h.chartsService.ConvertChartToTemplComponent(line)
	partials.BankChartCard(chartComponent).Render(r.Context(), w)
}
