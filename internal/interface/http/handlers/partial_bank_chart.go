package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// BankChartHandler renders the bank detail chart card: the month's daily
// balance line plus money in/out and transaction-count stats.
type BankChartHandler struct {
	accountQuery query.AccountQueryService
}

func NewBankChartHandler(
	aq query.AccountQueryService,
) *BankChartHandler {
	return &BankChartHandler{
		accountQuery: aq,
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

// Get renders the chart card; ?month=YYYY-MM selects the displayed month
// (defaulting to the current one). TODO: rethink, make it dynamic with
// month to date, last 6 months, last year, year to date and all
func (h *BankChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")

	now := time.Now()
	fotm := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	changeHist, err := h.accountQuery.GetAccountDailyBalanceAndStatsSince(r.Context(), id, fotm)
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

	theme := r.Header.Get("theme")

	histMap, _, max := charts.GenerateAccountsHistLine(changeHist.Data)
	line := charts.LineTimeAccount(
		histMap,
		0,
		float64(max),
		startDate,
		theme,
	)

	chartComponent := charts.ConvertChartToTemplComponent(line)
	partials.BankChartCard(chartComponent, changeHist.MoneyIn, changeHist.MoneyOut, changeHist.Transactions).
		Render(r.Context(), w)
}
