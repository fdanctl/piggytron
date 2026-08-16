package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// AccountsHistoryChartHandler renders the combined daily balance history
// line chart for all of the user's banks.
type AccountsHistoryChartHandler struct {
	accountQuery query.AccountQueryService
	ledgerQuery  query.LedgerQueryService
}

func NewAccountsHistoryChartHandler(
	aq query.AccountQueryService,
	lq query.LedgerQueryService,
) *AccountsHistoryChartHandler {
	return &AccountsHistoryChartHandler{
		accountQuery: aq,
		ledgerQuery:  lq,
	}
}

func (h *AccountsHistoryChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the banks' combined balance history from a selected
// predefind time period (all, 1y, ytd, 6m, mtd).
func (h *AccountsHistoryChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	theme := r.Header.Get("theme")
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "ytd" // default
	}

	// TODO cache it in redis
	d, err := h.ledgerQuery.GetFirstEntryDate(r.Context(), sessionInfo.UserID)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			d = time.Now()
		} else {
			httperror.SendError(
				w,
				r,
				errs.NewInternalAppError(err, "AccountsHistoryChartHandler.Get"),
			)
			return
		}
	}

	startDate := getStartPeriodDate(period, d)

	changeHist, err := h.accountQuery.GetAllDailyBalanceSince(
		r.Context(),
		sessionInfo.UserID,
		startDate,
	)
	if err != nil {
		if !errors.Is(err, query.ErrNoHistory) {
			httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
			return
		}
	}

	chart := components.NoData()

	if len(changeHist) > 0 {
		histMap, sortedKeys, _, min, max := charts.GenerateAccountsHistLineAndPieItems(
			changeHist,
		)
		chart = charts.ConvertChartToTemplComponent(charts.LineTime(
			histMap,
			&sortedKeys,
			min,
			max,
			startDate,
			theme,
		))
	}

	partials.AccountHistCard(chart, period).
		Render(r.Context(), w)
}
