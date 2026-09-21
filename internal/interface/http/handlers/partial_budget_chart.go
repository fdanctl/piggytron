package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// BudgetChartHandler renders the budget sankey diagram for a given month
// (path {month}), including unassigned carryover and
// leftover.
type BudgetChartHandler struct {
	categoryQuery query.CategoryQueryService
}

func NewBudgetChartHandler(
	cq query.CategoryQueryService,
) *BudgetChartHandler {
	return &BudgetChartHandler{
		categoryQuery: cq,
	}
}

func (h *BudgetChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the sankey for the requested month, honoring the optional
// ?unassign= carryover amount.
func (h *BudgetChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	month := r.PathValue("month")
	logger := middleware.LoggerFromContext(r.Context())
	logger.Debug(month)
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	punassign := r.URL.Query().Get("unassign")
	ponHold := r.URL.Query().Get("on-hold")

	var unassignCarry int
	if punassign != "" {
		unassignCarry, err = strconv.Atoi(punassign)
		if err != nil {
			httperror.SendError(w, r, err)
			return
		}
	}

	bm := budget.NewMonth(time.Now())
	if month != "" {
		bm, err = budget.ParseMonth(month)
		if err != nil {
			err := errs.NewAppError(
				errs.KindBadRequest,
				fmt.Sprintf("%s is not a valid month", month),
				fmt.Errorf("failed to parse month '%s': %w", month, err),
				"BudgetChartHandler.Get",
			)
			httperror.SendError(w, r, err)
			return
		}
	}

	categoryBudget, err := h.categoryQuery.GetCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		bm.Time(),
		time.Date(bm.Time().Year(), bm.Time().Month()+1, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		err := fmt.Errorf("error geting categories budget-spent: %w", err)
		httperror.SendError(w, r, err)
		return
	}

	onHold, err := strconv.Atoi(ponHold)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	theme := r.Header.Get("theme")
	component := charts.BudgetChart(unassignCarry, onHold, theme, categoryBudget)

	component.Render(r.Context(), w)
}
