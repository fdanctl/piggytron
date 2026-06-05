package handlers

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
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
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	)

	tWithCount, err := h.query.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*uint(page)-LIMIT,
	)
	if err != nil {
		logger.Error("error finding filtered transactions", "error", err, "filters", filters)
		http.Error(w, "internal error", http.StatusInternalServerError)
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

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		header := partials.AllTransactionsHeader([]components.BreadcrumbsLink{
			{Href: "", Name: "Transactions"},
			{Href: "", Name: "All"},
		},
			nil,
			uint8(filterCount),
			r.URL.RawQuery,
			tWithCount.Total,
		)
		if header.Render(ctx, w); err != nil {
			return err
		}

		return partials.TransactionsList(transactionsView, strings.Join(queries, "&"), hasMore).
			Render(ctx, w)
	})

	renderWithMainLayout(w, r, "Transactions", content)
}
