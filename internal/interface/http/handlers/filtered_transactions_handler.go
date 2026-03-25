package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type FilteredTransactionsHandler struct {
	query query.TransactionQueryService
}

func NewFilteredTransactionsHandler(
	q query.TransactionQueryService,
) *FilteredTransactionsHandler {
	return &FilteredTransactionsHandler{
		query: q,
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

	filters, err := query.NewTransactionFilters(types, accounts, cats, minAmount, maxAmount)
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

	v := r.Context().Value("user")
	if v == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	transactions, err := h.query.FindFiltered(
		r.Context(),
		sessionInfo.UserId,
		filters,
		LIMIT+1,
		LIMIT*uint(page)-LIMIT,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var hasMore bool
	if len(transactions) == LIMIT+1 {
		hasMore = true
		transactions = transactions[0 : len(transactions)-1]
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for i, v := range transactions {
			t := views.NewTransaction(v)
			var c templ.Component
			if i == len(transactions)-1 && hasMore {
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

	content.Render(r.Context(), w)
}
