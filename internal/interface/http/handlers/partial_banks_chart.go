package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
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
		httperror.SendError(w, r, errs.NewInternalAppError(err, "BanksChartsHandler.Get"))
		return
	}

	defaultPeriod := "ytd"
	startDate := getStartPeriodDate(defaultPeriod, d)

	changeHist, err := h.accountQuery.GetAllDailyBalanceSince(
		r.Context(),
		sessionInfo.UserID,
		startDate,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	histMap, sortedKeys, pieItems, min, max := charts.GenerateAccountsHistLineAndPieItems(
		changeHist,
	)
	line := charts.LineTime(
		histMap,
		&sortedKeys,
		min,
		max,
		startDate,
		theme,
	)

	c := charts.PieRadius(pieItems, "Assets", theme)
	pie := charts.ConvertChartToTemplComponent(c)

	templ.Join(
		pie,
		layouts.OOBWraper(
			"account-history-chart",
			"outerHTML",
			nil,
			partials.AccountHistCard(charts.ConvertChartToTemplComponent(line), defaultPeriod),
		),
	).Render(r.Context(), w)
}
