package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type DashboardHandler struct {
	ledgerQuery   query.LedgerQueryService
	accountQuery  query.AccountQueryService
	categoryQuery query.CategoryQueryService
}

func NewDashboardHandler(
	lq query.LedgerQueryService,
	aq query.AccountQueryService,
	cq query.CategoryQueryService,
) *DashboardHandler {
	return &DashboardHandler{
		ledgerQuery:   lq,
		accountQuery:  aq,
		categoryQuery: cq,
	}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	transactions, err := h.ledgerQuery.GetRecentEntries(
		r.Context(),
		sessionInfo.UserID,
		5,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts: %w", err))
		return
	}

	bm := budget.NewMonth(time.Now())
	categoryBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error geting category budget-spent: %w", err))
		return
	}

	pageView := views.NewDashboardPage(accounts, transactions, categoryBudgetSpent.Data)
	content := pages.Dashboard(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Dashboard"},
			},
			Options: nil,
		},
		pageView,
	)

	renderWithMainLayout(w, r, "Dashboard", content)
}
