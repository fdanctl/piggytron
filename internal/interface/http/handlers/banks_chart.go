package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/charts"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

type BanksChartsHandler struct {
	chartsService    *charts.Service
	transactionQuery query.TransactionQueryService
	accountQuery     query.AccountQueryService
}

func NewBanksChartsHandler(
	cs *charts.Service,
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
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error finding accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pieItems := h.chartsService.MakeAssetsPieItems(accounts, 5)
	pie := h.chartsService.PieRadius(pieItems)

	changeHist, err := h.accountQuery.GetBanksDailyChange(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error finding accounts history", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
