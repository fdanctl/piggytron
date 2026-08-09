package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/views/charts"
)

type BanksChartsHandler struct {
	accountQuery query.AccountQueryService
}

func NewBanksChartsHandler(
	aq query.AccountQueryService,
) *BanksChartsHandler {
	return &BanksChartsHandler{
		accountQuery: aq,
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

func (h *BanksChartsHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	// TODO: optimize could easily be just one query
	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts: %w", err))
		return
	}

	theme := r.Header.Get("theme")

	pieItems := charts.MakeAccountsPieItems(accounts)
	pie := components.NoData()
	if len(pieItems) > 0 {
		c := charts.PieRadius(pieItems, "Assets", theme)
		pie = charts.ConvertChartToTemplComponent(c)
	}

	changeHist, err := h.accountQuery.GetBanksDailyChange(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}
	histMap, min, max := charts.GenerateYearAccountsHistLine(changeHist)
	line := charts.LineTime(histMap, min, max, theme)

	templ.Join(
		pie,
		layouts.OOBWraper(
			"account-history-chart",
			"innerHTML",
			nil,
			charts.ConvertChartToTemplComponent(line),
		),
	).Render(r.Context(), w)
}
