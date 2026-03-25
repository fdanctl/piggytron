package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type AllTransactionsHandler struct {
	service *transactionapp.Service
}

func NewAllTransactionsHandler(s *transactionapp.Service) *AllTransactionsHandler {
	return &AllTransactionsHandler{
		service: s,
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

	filters, err := transaction.NewFilters(qtypes, qaccounts, qcats, qminAmount, qmaxAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterCount := len(qtypes) + len(qaccounts) + len(qcats)
	queries := []string{fmt.Sprintf("page=%d", 2)}
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

	transactions, resCount, hasMore, err := h.service.ReadFilteredWithCount(r.Context(), filters, 1)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	var transactionsView []views.Transaction
	for _, t := range transactions {
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
			resCount,
		)
		if header.Render(ctx, w); err != nil {
			return err
		}

		return partials.TransactionsList(transactionsView, r.URL.RawQuery, hasMore).
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
