package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type BudgetHandler struct {
	categoryQuery    query.CategoryQueryService
	transactionQuery query.TransactionQueryService
}

func NewBudgetHandler(
	cq query.CategoryQueryService,
	tq query.TransactionQueryService,
) *BudgetHandler {
	return &BudgetHandler{
		transactionQuery: tq,
		categoryQuery:    cq,
	}
}

func (h *BudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		logger.Error("error findingTransactions", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
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
		logger.Error("error geting category budget-spent", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	pageView := views.NewBudgetPageView(totalIncome, categoryBudgetSpent)

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Budget"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		err = partials.Budget(pageView).Render(ctx, w)
		return err
	})

	renderWithMainLayout(w, r, "Budget", content)
}
