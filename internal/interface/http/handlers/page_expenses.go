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

// ExpensesHandler renders the expenses report for a given month.
type ExpensesHandler struct {
	ledgerQuery query.LedgerQueryService
}

func NewExpensesHandler(
	lq query.LedgerQueryService,
) *ExpensesHandler {
	return &ExpensesHandler{
		ledgerQuery: lq,
	}
}

func (h *ExpensesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the month's categorized expenses, their total and the month
// selector; ?month=YYYY-MM selects the displayed month.
// TODO: WIP
func (h *ExpensesHandler) Get(w http.ResponseWriter, r *http.Request) {
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
				"ExpensesHandler.Post",
			)
			httperror.SendError(w, r, err)
			return
		}
	}

	nextMonthFirst := time.Date(
		bm.Time().Year(),
		bm.Time().Month()+1,
		1,
		0,
		0,
		0,
		0,
		bm.Time().Location(),
	)

	transactions, err := h.ledgerQuery.FindAllWithExpenseCategoryWithCount(
		r.Context(),
		sessionInfo.UserID,
		bm.Time(),
		nextMonthFirst,
		0,
		0,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error finding filtered transaction: %w", err))
		return
	}

	var totalSpent int
	var transactionsView []views.Transaction
	for _, t := range transactions.Data {
		totalSpent += t.Amount
		transactionsView = append(
			transactionsView,
			views.NewTransaction(t),
		)
	}

	// TODO cache it in redis
	d, err := h.ledgerQuery.GetFirstEntryDate(r.Context(), sessionInfo.UserID)
	if err != nil {
		if errors.Is(err, ledger.ErrNotFound) {
			d = time.Now()
		} else {
			httperror.SendError(w, r, errs.NewInternalAppError(err, "ExpensesPageHandler.Get"))
			return
		}
	}

	content := pages.Expenses(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Reports"},
				{Href: "", Name: "Expenses"},
			},
			Options: []views.BreadcrumbsLink{
				{Href: "/reports/expenses", Name: "Expenses"},
				{Href: "/reports/income", Name: "Income"},
			},
		},
		bm.Label(),
		d,
		transactionsView, transactions.Total,
		totalSpent,
	)

	renderWithMainLayout(w, r, "Expenses", content)
}
