package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// AccountChartHandler renders the daily balance line chart for one account.
type AccountChartHandler struct {
	accountQuery query.AccountQueryService
}

func NewAccountChartHandler(
	aq query.AccountQueryService,
) *AccountChartHandler {
	return &AccountChartHandler{
		accountQuery: aq,
	}
}

func (h *AccountChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

// Get renders the account's daily balance line from ?start=YYYY-MM-DD
// (defaulting to the start of the year). Used in goals page, start and max
// are used to render the chart.
func (h *AccountChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")

	q := r.URL.Query()
	d := q.Get("start")
	m := q.Get("max")
	logger.Debug(d)

	qmax, err := strconv.Atoi(m)
	if err != nil {
		qmax = 0
	}

	startDate, err := time.Parse(time.DateOnly, d)
	if err != nil {
		startDate = time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.Local)
	}

	changeHist, err := h.accountQuery.GetAccountDailyBalanceSince(
		r.Context(),
		id,
		startDate,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts history: %w", err))
		return
	}

	theme := r.Header.Get("theme")

	histMap, _, max := charts.GenerateAccountsHistLine(changeHist)
	line := charts.LineTimeAccount(
		histMap,
		0,
		math.Max(float64(qmax)/100, float64(max)),
		startDate,
		theme,
	)
	charts.ConvertChartToTemplComponent(line).Render(r.Context(), w)
}
