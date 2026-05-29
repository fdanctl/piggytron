package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type TransactionDetailsHandler struct {
	service query.TransactionQueryService
}

func NewTransactionDetailsHandler(s query.TransactionQueryService) *TransactionDetailsHandler {
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

func (h *TransactionDetailsHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	id := r.PathValue("id")
	t, err := h.service.FindByID(r.Context(), id)
	if err != nil {
		logger.Error("error finding transaction", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	tview := views.NewTransaction(*t)
	partials.TransactionDetails(tview).Render(r.Context(), w)
}
