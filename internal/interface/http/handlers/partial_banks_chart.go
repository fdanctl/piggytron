package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// BanksChartsHandler renders the banks page charts: the assets pie and the
// combined daily balance history line.
type BanksChartsHandler struct {
	accountQuery query.AccountQueryService
	ledgerQuery  query.LedgerQueryService
}

func NewBanksChartsHandler(
	aq query.AccountQueryService,
	lq query.LedgerQueryService,
) *BanksChartsHandler {
	return &BanksChartsHandler{
		accountQuery: aq,
		ledgerQuery:  lq,
	}
}

func (h *BanksChartsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the assets pie and the accounts balance chart (out-of-band).
func (h *BanksChartsHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	theme := r.Header.Get("theme")

	d, err := h.ledgerQuery.GetFirstEntryDate(r.Context(), sessionInfo.UserID)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			d = time.Now()
		} else {
			httperror.SendError(
				w,
				r,
				errs.NewInternalAppError(err, "BanksChartsHandler.Get"),
			)
			return
		}
	}

	defaultPeriod := "ytd"
	startDate := getStartPeriodDate(defaultPeriod, d)

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

	line := components.NoData()
	pie := components.NoData()

	if len(changeHist) > 0 {
		histMap, sortedKeys, pieItems, min, max := charts.GenerateAccountsHistLineAndPieItems(
			changeHist,
		)
		line = charts.ConvertChartToTemplComponent(charts.LineTime(
			histMap,
			&sortedKeys,
			min,
			max,
			startDate,
			theme,
		))
		pie = charts.ConvertChartToTemplComponent(charts.PieRadius(pieItems, "Assets", theme))
	}

	templ.Join(
		pie,
		layouts.OOBWraper(
			"account-history-chart",
			"outerHTML",
			nil,
			partials.AccountHistCard(line, defaultPeriod),
		),
	).Render(r.Context(), w)
}
