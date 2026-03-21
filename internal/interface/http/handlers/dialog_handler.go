package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	bankapp "github.com/fdanctl/piggytron/internal/application/bank"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategoryapp "github.com/fdanctl/piggytron/internal/application/income_category"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type DialogHandler struct {
	incomeCatService   *incomecategoryapp.Service
	expenseCatService  *expensecategoryapp.Service
	transactionService *transactionapp.Service
	bankService        *bankapp.Service
}

func NewDialogHandler(
	es *expensecategoryapp.Service,
	is *incomecategoryapp.Service,
	ts *transactionapp.Service,
	bs *bankapp.Service,
) *DialogHandler {
	return &DialogHandler{
		incomeCatService:   is,
		expenseCatService:  es,
		transactionService: ts,
		bankService:        bs,
	}
}

func (h *DialogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		dialog := r.PathValue("dialog")

		switch dialog {
		case "transaction-filters":
			h.GetDialogFilters(w, r)

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

	content := partials.TransactionsFilters(categoryOptions, accountOptions, r.URL.Query())
	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("sheet", nil).Render(ctx, w)
}
