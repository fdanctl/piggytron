package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/domain/ledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

// BudgetPageHandler renders the budget page for a given month: budgets per
// category, spent values and the month's net/balance.
type BudgetPageHandler struct {
	categoryQuery    query.CategoryQueryService
	transactionQuery query.LedgerQueryService
}

func NewBudgetPageHandler(
	cq query.CategoryQueryService,
	tq query.LedgerQueryService,
) *BudgetPageHandler {
	return &BudgetPageHandler{
		transactionQuery: tq,
		categoryQuery:    cq,
	}
}

func (h *BudgetPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the budget page; the optional ?month=YYYY-MM query selects the
// displayed month.
func (h *BudgetPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	month := r.URL.Query().Get("month")

	bm := budget.NewMonth(time.Now())
	if month != "" {
		bm, err = budget.ParseMonth(month)
		if err != nil {
			err := errs.NewAppError(
				errs.KindBadRequest,
				fmt.Sprintf("%s is not a valid month", month),
				fmt.Errorf("failed to parse month '%s': %w", month, err),
				"BudgetHandler.Post",
			)
			httperror.SendError(w, r, err)
			return
		}
	}

	categoryBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error geting category budget-spent: %w", err))
		return
	}

	pageView := views.NewBudgetPageView(
		bm,
		categoryBudgetSpent.MonthNet,
		categoryBudgetSpent.Balance,
		categoryBudgetSpent.Data,
	)
	d, err := h.transactionQuery.GetFirstEntryDate(r.Context(), sessionInfo.UserID)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			d = time.Now()
		} else {
			httperror.SendError(w, r, errs.NewInternalAppError(err, "BudgetPageHandler.Get"))
			return
		}
	}

	content := pages.Budget(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Budget"},
			},
			Options: nil,
		},
		d,
		pageView,
	)

	renderWithMainLayout(w, r, "Budget", content)
}
