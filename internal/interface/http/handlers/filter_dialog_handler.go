package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	categorynameapp "github.com/fdanctl/piggytron/internal/application/category_name"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type FilterDialogHandler struct {
	categorynameService *categorynameapp.Service
	transactionService  *transactionapp.Service
	accountService      *accountapp.Service
}

func NewFilterDialogHandler(
	cs *categorynameapp.Service,
	ts *transactionapp.Service,
	as *accountapp.Service,
) *FilterDialogHandler {
	fmt.Println(cs)
	return &FilterDialogHandler{
		categorynameService: cs,
		transactionService:  ts,
		accountService:      as,
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
	c, err := h.categorynameService.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var categoryOptions []partials.FilterOption
	for _, v := range c {
		categoryOptions = append(
			categoryOptions,
			partials.FilterOption{Label: v.Name, Value: v.Id},
		)
	}

	account, err := h.accountService.ReadAllByUser(r.Context())
	if err != nil {
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

	filters, err := transaction.NewFilters(types, accounts, cats, minAmount, maxAmount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	resCount, err := h.transactionService.CountFilteredResults(r.Context(), filters)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	includedAcc, err := h.accountService.ReadIdNamesIncludes(r.Context(), accounts)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	includedCats, err := h.categorynameService.GetCategoriesIdIncludes(r.Context(), cats)
	if err != nil {
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
	r.ParseForm()
	q := r.Form
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

	resCount, err := h.transactionService.CountFilteredResults(r.Context(), filters)
	if err != nil {
		fmt.Println(err)
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
		"hx-get":      "/partials/dialog/transaction-filters?" + strings.Join(queries[1:], "&"),
		"hx-target":   "#dialog-root",
	}).Render(r.Context(), w)
}
