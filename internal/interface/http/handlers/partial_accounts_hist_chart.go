package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// AccountsHistoryChartHandler renders the combined daily balance history
// line chart for all of the user's banks.
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

// Get renders the banks' combined balance history from January 1st of the
// current year. TODO: rethink, maybe it makes more sense to be one year range
// better yet make it dynamic
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

	theme := r.Header.Get("theme")

	changeHist, err := h.accountQuery.GetAllDailyBalanceSince(
		r.Context(),
		sessionInfo.UserID,
		time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local),
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	histMap, sortedKeys, _, min, max := charts.GenerateAccountsHistLineAndPieItems(
		changeHist,
	)
	line := charts.LineTime(
		histMap,
		&sortedKeys,
		min,
		max,
		time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local),
		theme,
	)

	charts.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
