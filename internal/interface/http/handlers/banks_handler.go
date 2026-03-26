package handlers

import (
	"fmt"
	"net/http"

	accountapp "github.com/fdanctl/piggytron/internal/application/account"
)

type BanksHandler struct {
	service *accountapp.Service
}

func NewBanksHandler(s *accountapp.Service) *BanksHandler {
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
		h.GetWithID(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BanksHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := sessionInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banks, err := h.service.ReadAllBanksByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	for _, b := range banks {
		fmt.Fprintf(w, "<a href=\"/banks/%s\">%s</p>", b.ID(), b.Name())
	}
}

func (h *BanksHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	bank, err := h.service.ReadOneByID(r.Context(), r.PathValue("id"))
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<p>%s</p>", bank.Name())
	fmt.Fprintf(w, "<p>%s</p>", bank.Currency())
}
