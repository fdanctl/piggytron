package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
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
	q := r.URL.Query()
	qtypes := q["types"]
	qaccounts := q["accounts"]
	qcats := q["categories"]
	qminAmount := q.Get("minamount")
	qmaxAmount := q.Get("maxamount")

	filters, err := query.NewTransactionFilters(qtypes, qaccounts, qcats, qminAmount, qmaxAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterCount := len(qtypes) + len(qaccounts) + len(qcats)
	var page uint = 1
	queries := []string{fmt.Sprintf("page=%d", page+1)}
	if len(qtypes) > 0 {
		queries = append(queries, "types="+strings.Join(qtypes, "&types="))
	}
	if len(qaccounts) > 0 {
		queries = append(queries, "accounts="+strings.Join(qaccounts, "&accounts="))
	}
	if len(qcats) > 0 {
		queries = append(queries, "categories="+strings.Join(qcats, "&categories="))
	}
	if qminAmount != "" {
		queries = append(queries, "minamount="+qminAmount)
		filterCount++
	}
	if qmaxAmount != "" {
		queries = append(queries, "maxmount="+qminAmount)
		filterCount++
	}

	sessionInfo, err := sessionInfoFormCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tWithCount, err := h.query.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*page-LIMIT,
	)
	if err != nil {
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

	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
		io.WriteString(w, "<title>Transactions</title>")
		return
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Transactions").Render(ctx, w)
}
