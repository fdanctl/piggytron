package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// TransactionDetailsHandler renders the read-only details dialog for one
// ledger entry.
type TransactionDetailsHandler struct {
	service query.LedgerQueryService
}

func NewTransactionDetailsHandler(s query.LedgerQueryService) *TransactionDetailsHandler {
	return &TransactionDetailsHandler{
		service: s,
	}
}

func (h *TransactionDetailsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the details dialog with the entry's action buttons.
func (h *TransactionDetailsHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	t, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	tview := views.NewTransaction(*t)
	components.DialogWrapper(
		"dialog--top-right",
		components.DialogHeader("", "Transaction Details", nil),
		partials.TransactionDetails(tview),
		partials.TransactionDetailsButtons(tview),
		nil,
	).
		Render(r.Context(), w)
}
