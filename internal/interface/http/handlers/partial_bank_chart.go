package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/account"
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

// Get renders the chart card. Shows account balance chart from a
// selected predefind time period (all, 1y, ytd, 6m, mtd), and
// the money in, money out and transaction count from that period.
func (h *BankChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")
	theme := r.Header.Get("theme")

	period := r.URL.Query().Get("period")
	if period == "" {
		period = "mtd" // default
	}

	qclosedAt := r.URL.Query().Get("closedAt")
	logger.Debug("closed", "closed", qclosedAt)
	now := time.Now()
	end := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	var err error = nil
	if qclosedAt != "" {
		end, err = time.Parse(time.DateOnly, qclosedAt)
		if err != nil {
			httperror.SendError(w, r, err)
		}
		theme += "-closed"
		period = "all"
	}

	logger.Debug("closed", "closed", end)

	// TODO cache it in redis
	d, err := h.accountQuery.GetAccountFirstEntryDate(r.Context(), id)
	if err != nil {
		if errors.Is(err, account.ErrNotFound) {
			d = end
		} else {
			httperror.SendError(w, r, fmt.Errorf("failed to accounts first ledger entry: %w", err))
			return
		}
	}

	startDate := getStartPeriodDate(period, d)

	changeHist, err := h.accountQuery.GetAccountDailyBalanceAndStatsSince(
		r.Context(),
		id,
		startDate,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	histMap, _, max := charts.GenerateAccountsHistLine(changeHist.Data)
	line := charts.LineTimeAccount(
		histMap,
		0,
		float64(max),
		startDate,
		end,
		theme,
	)

	chartComponent := charts.ConvertChartToTemplComponent(line)
	partials.BankChartCard(id, chartComponent, changeHist.MoneyIn, changeHist.MoneyOut, changeHist.Transactions, period, qclosedAt == "").
		Render(r.Context(), w)
}
