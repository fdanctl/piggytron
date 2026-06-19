package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type BudgetPageHandler struct {
	categoryQuery    query.CategoryQueryService
	transactionQuery query.TransactionQueryService
}

func NewBudgetPageHandler(
	cq query.CategoryQueryService,
	tq query.TransactionQueryService,
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

func (h *BudgetPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	// TODO url month query (ex. 202605)
	now := time.Now()

	filters := query.NewTransactionFilters([]string{"income"}, nil, nil, "", "", "", "")
	minD := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	maxD := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	filters.MinDate = &minD
	filters.MaxDate = &maxD

	incomes, err := h.transactionQuery.FindFiltered(
		r.Context(),
		sessionInfo.UserID,
		filters,
		0,
		0,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find filtered transactions: %w", err))
		return
	}

	var totalIncome int
	for _, v := range incomes {
		totalIncome += v.Amount
	}
	logger.Debug("total income: " + fmt.Sprint(totalIncome))

	categoryBudgetSpent, err := h.categoryQuery.GetExpenseCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		minD,
		maxD,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error geting category budget-spent: %w", err))
		return
	}

	pageView := views.NewBudgetPageView(totalIncome, categoryBudgetSpent)

	content := pages.Budget(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Budget"},
			},
			Options: nil,
		},
		pageView,
	)

	renderWithMainLayout(w, r, "Budget", content)
}
