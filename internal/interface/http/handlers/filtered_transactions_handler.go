package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type FilteredTransactionsHandler struct {
	service *transactionapp.Service
}

func NewFilteredTransactionsHandler(s *transactionapp.Service) *FilteredTransactionsHandler {
	return &FilteredTransactionsHandler{
		service: s,
	}
}

func (h *FilteredTransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FilteredTransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	pageq := r.URL.Query().Get("page")
	if pageq == "" {
		pageq = "1"
	}
	page, err := strconv.Atoi(pageq)
	if err != nil {
		http.Error(w, "page must must be a number", http.StatusBadRequest)
		return
	}
	if page < 1 {
		http.Error(w, "page must must be positive", http.StatusBadRequest)
		return
	}

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
	queries := []string{fmt.Sprintf("page=%d", page+1)}
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
	fmt.Printf("queryStr: %v\n", strings.Join(queries, "&"))

	transactions, err := h.service.ReadWithFilters(r.Context(), filters, uint(page))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for i, v := range transactions {
			t := views.NewTransaction(v.ID(), v.Description(), v.Ttype(), v.Amount(), v.Date())
			var c templ.Component
			if i == len(transactions)-1 {
				c = partials.TransactionItem(t, templ.Attributes{
					"style": fmt.Sprintf("animation-delay: %dms", i*30),
					"hx-get": fmt.Sprintf(
						"/partials/transactions?%s",
						strings.Join(queries, "&"),
					),
					"hx-trigger":   "intersect once",
					"hx-swap":      "afterend",
					"hx-indicator": "#infinite-indicator",
				})
			} else {
				c = partials.TransactionItem(t, templ.Attributes{
					"style": fmt.Sprintf("animation-delay: %dms", i*30),
				})
			}
			err := c.Render(ctx, w)
			if err != nil {
				return err
			}
		}
		return nil
	})

	w.Header().Set("HX-Push-Url", "?"+strings.Join(queries[1:], "&"))
	content.Render(r.Context(), w)
	components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
		"hx-swap-oob": "outerHTML",
		"id":          "filter-btn",
		"hx-get":      "/partials/dialog/transaction-filters?" + r.URL.RawQuery,
		"hx-target":   "#dialog-root",
	}).Render(r.Context(), w)
}
