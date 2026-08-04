package handlers

import (
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/views/charts"
)

type AccountsHistoryChartHandler struct {
	accountQuery query.AccountQueryService
}

func NewAccountsHistoryChartHandler(
	aq query.AccountQueryService,
) *AccountsHistoryChartHandler {
	return &AccountsHistoryChartHandler{
		accountQuery: aq,
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
	histMap, min, max := charts.GenerateYearAccountsHistLine(changeHist)
	line := charts.LineTime(histMap, min, max)

	charts.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
