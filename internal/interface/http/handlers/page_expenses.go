package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/web/templates/pages"
)

type ExpensesHandler struct{}

func (h *ExpensesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExpensesHandler) Get(w http.ResponseWriter, r *http.Request) {
	content := pages.Test()
	renderWithMainLayout(w, r, "Expenses", content)
}
