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

type GoalContributionsHandler struct {
	query query.TransactionQueryService
}

func NewGoalContributionsHandler(
	q query.TransactionQueryService,
) *GoalContributionsHandler {
	return &GoalContributionsHandler{
		query: q,
	}
}

func (h *GoalContributionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GoalContributionsHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	accounts := q["accounts"]
	name := q.Get("name")

	filters := query.NewTransactionFilters(
		nil,
		accounts,
		nil,
		"",
		"",
		"",
		"",
	)

	queries := []string{fmt.Sprintf("page=%d", page+1)}
	if len(accounts) > 0 {
		queries = append(queries, "accounts="+strings.Join(accounts, "&accounts="))
	}
	queries = append(queries, "name="+name)

	transactions, err := h.query.FindFiltered(
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
	if len(transactions) == LIMIT+1 {
		hasMore = true
		transactions = transactions[0 : len(transactions)-1]
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		for i, v := range transactions {
			t := views.NewAccountTransaction(v, name)
			var c templ.Component
			if i == len(transactions)-1 && hasMore {
				c = partials.ContributionItem(t, templ.Attributes{
					"style": fmt.Sprintf("animation-delay: %dms", i*30),
					"hx-get": fmt.Sprintf(
						"/partials/contributions?%s",
						strings.Join(queries, "&"),
					),
					"hx-trigger":   "intersect once",
					"hx-swap":      "afterend",
					"hx-indicator": "#infinite-indicator",
				})
			} else {
				c = partials.ContributionItem(t, templ.Attributes{
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
