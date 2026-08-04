package handlers

import (
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
)

type AccountsHistoryChartHandler struct {
	chartsService *appcharts.Service
	accountQuery  query.AccountQueryService
}

func NewAccountsHistoryChartHandler(
	cs *appcharts.Service,
	aq query.AccountQueryService,
) *AccountsHistoryChartHandler {
	return &AccountsHistoryChartHandler{
		chartsService: cs,
		accountQuery:  aq,
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

func (h *AccountsHistoryChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	period := r.URL.Query().Get("period")
	if period == "" {
		period = "ytd"
	}

	changeHist, err := h.accountQuery.GetBanksDailyChange(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}
	histMap, min, max := h.chartsService.GenerateYearAccountsHistLine(changeHist)
	line := h.chartsService.LineTime(histMap, min, max)

	h.chartsService.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
