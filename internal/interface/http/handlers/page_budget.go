package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

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

func (h *BudgetPageHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
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
	logger.Debug(bm.String())

	filters := query.NewLedgerFilters([]string{"income"}, nil, nil, "", "", "", "")
	minD := bm.Time()
	maxD := time.Date(bm.Time().Year(), bm.Time().Month()+1, 1, 0, 0, 0, 0, time.UTC)
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
		httperror.SendError(w, r, fmt.Errorf("failed to find filtered ledger entries: %w", err))
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

	pageView := views.NewBudgetPageView(bm, totalIncome, categoryBudgetSpent)

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
