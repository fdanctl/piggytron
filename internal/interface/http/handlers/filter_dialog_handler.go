package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type FilterDialogHandler struct {
	categoryQueryService query.CategoryQueryService
	accountService       *accountapp.Service
	tQueryService        query.TransactionQueryService
	accQueryService      query.AccountQueryService
}

func NewFilterDialogHandler(
	cs query.CategoryQueryService,
	as *accountapp.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *FilterDialogHandler {
	return &FilterDialogHandler{
		categoryQueryService: cs,
		accountService:       as,
		tQueryService:        tq,
		accQueryService:      aq,
	}
}

func (h *FilterDialogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FilterDialogHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := h.categoryQueryService.FindAllCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error finding all Categories", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var categoryOptions []partials.FilterOption
	for _, v := range c {
		categoryOptions = append(
			categoryOptions,
			partials.FilterOption{Label: v.Name, Value: v.ID},
		)
	}

	account, err := h.accountService.ReadAllByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error finding all accounts", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var accountOptions []partials.FilterOption
	for _, v := range account {
		accountOptions = append(
			accountOptions,
			partials.FilterOption{Label: v.Name(), Value: string(v.ID())},
		)
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

	resCount, err := h.tQueryService.CountFilteredResults(
		r.Context(), sessionInfo.UserID, filters,
	)
	if err != nil {
		logger.Error("counting filters transactions results", "error", err, "filters", filters)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	includedAcc, err := h.accQueryService.FindIDNamesIncludes(r.Context(), accounts)
	if err != nil {
		logger.Error("find accounts with ids array", "error", err, "accounts_ids", accounts)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	includedCats, err := h.categoryQueryService.FindCategoriesIDIncludes(r.Context(), cats)
	if err != nil {
		logger.Error("find category with ids array", "error", err, "categories_ids", cats)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	content := partials.TransactionsFilters(
		accountOptions,
		categoryOptions,
		includedAcc,
		includedCats,
		r.URL.Query(),
		resCount,
	)
	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("sheet", nil).Render(ctx, w)
}

func (h *FilterDialogHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	r.ParseForm()
	q := r.Form
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

	resCount, err := h.tQueryService.CountFilteredResults(
		r.Context(), sessionInfo.UserID, filters,
	)
	if err != nil {
		logger.Error("counting filters transactions results", "error", err, "filters", filters)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("HX-Push-Url", "?"+strings.Join(queries[1:], "&"))
	w.Header().Set("HX-Trigger", "refetch-transactions")

	if filterCount > 0 {
		components.Button(
			"Reset",
			"w-full",
			components.BtnOutline,
			components.BtnMedium,
			templ.Attributes{"type": "button", "onclick": "resetTransactionFiltersForm()"},
		).Render(r.Context(), w)
	}
	components.Button(
		fmt.Sprintf("Show %d results", resCount),
		"w-full",
		components.BtnPrimary,
		components.BtnMedium,
		templ.Attributes{"type": "button", "onclick": "closeLastDialog()"},
	).Render(r.Context(), w)

	templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p id=\"filter-result-count\" hx-swap-oob=\"innerHTML\">")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%d results", resCount)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</p>")
		return err
	}).Render(r.Context(), w)
	components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
		"style":       "height: 24px;",
		"hx-swap-oob": "outerHTML",
		"id":          "filter-btn",
		"hx-get":      "/partials/transaction-filters?" + strings.Join(queries[1:], "&"),
		"hx-target":   "#dialog-root",
	}).Render(r.Context(), w)
}
