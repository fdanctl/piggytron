package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategoryapp "github.com/fdanctl/piggytron/internal/application/income_category"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type DialogHandler struct {
	incomeCatService   *incomecategoryapp.Service
	expenseCatService  *expensecategoryapp.Service
	transactionService *transactionapp.Service
	bankService        *accountapp.Service
}

func NewDialogHandler(
	es *expensecategoryapp.Service,
	is *incomecategoryapp.Service,
	ts *transactionapp.Service,
	as *accountapp.Service,
) *DialogHandler {
	return &DialogHandler{
		incomeCatService:   is,
		expenseCatService:  es,
		transactionService: ts,
		bankService:        as,
	}
}

func (h *DialogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	dialog := r.PathValue("dialog")
	switch r.Method {
	case http.MethodGet:

		switch dialog {
		case "transaction-filters":
			h.GetDialogFilters(w, r)

		default:
			http.NotFound(w, r)

		}

	case http.MethodPost:
		switch dialog {
		case "transaction-filters":
			h.PostDialogFilters(w, r)

		default:
			http.NotFound(w, r)

		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DialogHandler) GetDialogFilters(w http.ResponseWriter, r *http.Request) {
	ic, err := h.incomeCatService.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	ec, err := h.expenseCatService.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var categoryOptions []partials.FilterOption
	for _, v := range ic {
		categoryOptions = append(
			categoryOptions,
			partials.FilterOption{Label: v.Name(), Value: string(v.ID())},
		)
	}
	for _, v := range ec {
		categoryOptions = append(
			categoryOptions,
			partials.FilterOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	banks, err := h.bankService.ReadAllByUser(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var accountOptions []partials.FilterOption
	for _, v := range banks {
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

	content := partials.TransactionsFilters(
		accountOptions,
		categoryOptions,
		r.URL.Query(),
		resCount,
	)
	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("sheet", nil).Render(ctx, w)
}

func (h *DialogHandler) PostDialogFilters(w http.ResponseWriter, r *http.Request) {
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

	fmt.Fprintln(w, fmt.Sprintf("Show %d results", resCount))
	components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
		"hx-swap-oob": "outerHTML",
		"id":          "filter-btn",
		"hx-get":      "/partials/dialog/transaction-filters?" + strings.Join(queries[1:], "&"),
		"hx-target":   "#dialog-root",
	}).Render(r.Context(), w)
}
