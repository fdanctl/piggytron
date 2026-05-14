package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/a-h/templ"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type BudgetHandler struct {
	expenseCatService *expensecategoryapp.Service
	transactionQuery  query.TransactionQueryService
}

func NewBudgetHandler(
	es *expensecategoryapp.Service,
	q query.TransactionQueryService,
) *BudgetHandler {
	return &BudgetHandler{
		expenseCatService: es,
		transactionQuery:  q,
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

	catExpenses, err := h.transactionQuery.GetExpensesByCategoryBetweenDates(
		r.Context(),
		sessionInfo.UserID,
		minD,
		maxD,
	)
	if err != nil {
		logger.Error("error getings expenses by category", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	logger.Debug("total expense: " + fmt.Sprint(catExpenses.Total))
	logger.Debug("total expense: " + fmt.Sprint(catExpenses.Data))

	ec, err := h.expenseCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	var ecView []views.ExpenseCategory
	for _, v := range ec {
		ecView = append(ecView, views.NewExpenseCategory(v))
	}

	leftToSpend := totalIncome - catExpenses.Total

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Budget"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		err = partials.Budget(totalIncome, leftToSpend, ecView).Render(ctx, w)
		return err
	})

	renderWithMainLayout(w, r, "Budget", content)
}
