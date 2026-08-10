package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/web/templates/pages"
)

// TestHandler renders the test page (development-only template sandbox).
type TestHandler struct{}

func (h *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the test page.
func (h *TestHandler) Get(w http.ResponseWriter, r *http.Request) {
	content := pages.Test()
	renderWithMainLayout(w, r, "Expenses", content)
}
