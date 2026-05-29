package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
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
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
	minDate := q.Get("mindate")
	maxDate := q.Get("maxdate")

	filters := query.NewTransactionFilters(
		types,
		accounts,
		cats,
		minAmount,
		maxAmount,
		minDate,
		maxDate,
	)

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

	transactions, err := h.query.FindFilteredWithCount(
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
	if len(transactions.Data) == LIMIT+1 {
		hasMore = true
		transactions.Data = transactions.Data[0 : len(transactions.Data)-1]
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for i, v := range transactions.Data {
			t := views.NewTransaction(v)
			var c templ.Component
			if i == len(transactions.Data)-1 && hasMore {
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

	obb := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p id=\"filter-result-count\" hx-swap-oob=\"innerHTML\">")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%d results", transactions.Total)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</p>")
		return err
	})
	templ.Join(content, obb).Render(r.Context(), w)
}
