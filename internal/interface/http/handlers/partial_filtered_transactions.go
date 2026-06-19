package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
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
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	pageq := r.URL.Query().Get("page")
	if pageq == "" {
		pageq = "1"
	}
	page, err := strconv.Atoi(pageq)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	if page < 1 {
		httperror.SendError(w, r, err)
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

	_, queries := queryStrFromFiltersWithCount(page+1, types, accounts, cats, minAmount, maxAmount)

	transactions, err := h.query.FindFilteredWithCount(
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
