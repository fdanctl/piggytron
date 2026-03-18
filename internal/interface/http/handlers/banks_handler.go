package handlers

import (
	"fmt"
	"net/http"

	bankapp "github.com/fdanctl/piggytron/internal/application/bank"
)

type BanksHandler struct {
	service *bankapp.Service
}

func NewBanksHandler(s *bankapp.Service) *BanksHandler {
	return &BanksHandler{service: s}
}

func (h *BanksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if id == "" {
			h.Get(w, r)
			return
		}
		h.GetId(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BanksHandler) Get(w http.ResponseWriter, r *http.Request) {
	banks, err := h.service.ReadAllByUser(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	for _, b := range banks {
		fmt.Fprintf(w, "<a href=\"/banks/%s\">%s</p>", b.ID(), b.Name())
	}
}

func (h *BanksHandler) GetId(w http.ResponseWriter, r *http.Request) {
	bank, err := h.service.ReadOneById(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<p>%s</p>", bank.Name())
	fmt.Fprintf(w, "<p>%s</p>", bank.Currency())
}
