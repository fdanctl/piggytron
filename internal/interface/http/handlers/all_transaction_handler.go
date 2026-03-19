package handlers

import (
	"fmt"
	"net/http"

	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
)

type AllTransactionsHandler struct {
	service *transactionapp.Service
}

func NewAllTransactionsHandler(s *transactionapp.Service) *AllTransactionsHandler {
	return &AllTransactionsHandler{
		service: s,
	}
}

func (h *AllTransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AllTransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO infinite scroll (maybe 50 each time)
	transactions, err := h.service.RealAllByUser(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	for _, b := range transactions {
		fmt.Fprintf(w, "<a href=\"/transactions/%s\">%s</p>", b.ID(), b.Description())
	}
}
