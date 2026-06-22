package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

type BanksChartsHandler struct {
	chartsService    *appcharts.Service
	transactionQuery query.TransactionQueryService
	accountQuery     query.AccountQueryService
}

func NewBanksChartsHandler(
	cs *appcharts.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *BanksChartsHandler {
	return &BanksChartsHandler{
		chartsService:    cs,
		transactionQuery: tq,
		accountQuery:     aq,
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

	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts: %w", err))
		return
	}

	pieItems := h.chartsService.MakeAssetsPieItems(accounts, 5)
	pie := h.chartsService.PieRadius(pieItems)

	changeHist, err := h.accountQuery.GetBanksDailyChange(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}
	histMap, min, max := h.chartsService.GenerateYearAccountsHistLine(changeHist)
	line := h.chartsService.LineTime(histMap, min, max)

	templ.Join(
		h.chartsService.ConvertChartToTemplComponent(pie),
		layouts.OOBWraper(
			"account-history-chart",
			"innerHTML",
			nil,
			h.chartsService.ConvertChartToTemplComponent(line),
		),
	).Render(r.Context(), w)
}
