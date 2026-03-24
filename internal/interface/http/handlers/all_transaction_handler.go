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
	types := q["types"]
	accounts := q["accounts"]
	cats := q["categories"]
	minAmount := q.Get("minamount")
	maxAmount := q.Get("maxamount")

	filters, err := transaction.NewFilters(types, accounts, cats, minAmount, maxAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	filterCount := len(types) + len(accounts) + len(cats)
	queries := []string{fmt.Sprintf("page=%d", 2)}
	if len(types) > 0 {
		queries = append(queries, "types="+strings.Join(types, "&types="))
	}
	if len(accounts) > 0 {
		queries = append(queries, "accounts="+strings.Join(accounts, "&accounts="))
	}
	if len(cats) > 0 {
		queries = append(queries, "categories="+strings.Join(cats, "&categories="))
	}
	if minAmount != "" {
		queries = append(queries, "minamount="+minAmount)
		filterCount++
	}
	if maxAmount != "" {
		queries = append(queries, "maxmount="+minAmount)
		filterCount++
	}

	transactions, hasMore, err := h.service.ReadWithFilters(r.Context(), filters, 1)
	if err != nil {
		fmt.Println(err)
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
		if _, err := io.WriteString(w, "<div class=\"flex justify-between\">"); err != nil {
			return err
		}
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Transactions"},
			{Href: "", Name: "All"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}
		if err := components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
			"id":        "filter-btn",
			"hx-get":    "/partials/dialog/transaction-filters?" + r.URL.RawQuery,
			"hx-target": "#dialog-root",
			"hx-swap":   "beforeend",
		}).Render(ctx, w); err != nil {
			return err
		}
		if _, err := io.WriteString(w, "</div>"); err != nil {
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
