package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type AllTransactionsHandler struct {
	query query.TransactionQueryService
}

func NewAllTransactionsHandler(
	q query.TransactionQueryService,
) *AllTransactionsHandler {
	return &AllTransactionsHandler{
		query: q,
	}
}

func (h *AllTransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AllTransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	q := r.URL.Query()
	qtypes := q["types"]
	qaccounts := q["accounts"]
	qcats := q["categories"]
	qminAmount := q.Get("minamount")
	qmaxAmount := q.Get("maxamount")
	qminDate := q.Get("mindate")
	qmaxDate := q.Get("maxdate")

	filters := query.NewTransactionFilters(
		qtypes,
		qaccounts,
		qcats,
		qminAmount,
		qmaxAmount,
		qminDate,
		qmaxDate,
	)

	page := 1
	filterCount, queries := queryStrFromFiltersWithCount(
		page+1,
		qtypes,
		qaccounts,
		qcats,
		qminAmount,
		qmaxAmount,
		qminDate,
		qmaxDate,
	)

	tWithCount, err := h.query.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*uint(page)-LIMIT,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find filtered transactions: %w", err))
		return
	}
	var hasMore bool
	if len(tWithCount.Data) == LIMIT+1 {
		hasMore = true
		tWithCount.Data = tWithCount.Data[0 : len(tWithCount.Data)-1]
	}

	var transactionsView []views.Transaction
	for _, t := range tWithCount.Data {
		transactionsView = append(
			transactionsView,
			views.NewTransaction(t),
		)
	}

	content := pages.AllTransactions(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Transactions"},
				{Href: "", Name: "All"},
			},
			Options: nil,
		},
		uint8(filterCount),
		r.URL.RawQuery,
		tWithCount.Total,
		transactionsView, strings.Join(queries, "&"), hasMore,
	)

	renderWithMainLayout(w, r, "Transactions", content)
}
