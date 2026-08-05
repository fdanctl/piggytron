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
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type IncomeHandler struct {
	ledgerQuery query.LedgerQueryService
}

func NewIncomeHandler(
	lq query.LedgerQueryService,
) *IncomeHandler {
	return &IncomeHandler{
		ledgerQuery: lq,
	}
}

func (h *IncomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *IncomeHandler) Get(w http.ResponseWriter, r *http.Request) {
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
				"IncomeHandler.Post",
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
	monthLast := time.Date(
		nextMonthFirst.Year(),
		nextMonthFirst.Month(),
		nextMonthFirst.Day()-1,
		0,
		0,
		0,
		0,
		nextMonthFirst.Location(),
	)

	// TODO cache it in redis
	d, err := h.ledgerQuery.GetFirstEntryDate(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, errs.NewInternalAppError(err, "IncomeHandler.Get"))
		return
	}

	filters := query.NewLedgerFilters(
		[]string{"income"},
		nil,
		nil,
		"",
		"",
		strconv.Itoa(int(bm.Time().Unix())),
		strconv.Itoa(int(monthLast.Unix())),
	)

	transactions, err := h.ledgerQuery.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		0,
		0,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error finding filtered transaction: %w", err))
		return
	}

	var total int
	sources := make(map[string]interface{}, 0)
	var transactionsView []views.Transaction
	for _, t := range transactions.Data {
		_, ok := sources[t.ID]
		if !ok {
			sources[t.ID] = nil
		}
		total += t.Amount
		transactionsView = append(
			transactionsView,
			views.NewTransaction(t),
		)
	}

	content := pages.Income(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Reports"},
				{Href: "", Name: "Income"},
			},
			Options: []views.BreadcrumbsLink{
				{Href: "/reports/expenses", Name: "Expenses"},
				{Href: "/reports/income", Name: "Income"},
			},
		},
		bm.Label(),
		d,
		transactionsView, transactions.Total,
		total,
		len(sources),
	)

	renderWithMainLayout(w, r, "Income", content)
}
