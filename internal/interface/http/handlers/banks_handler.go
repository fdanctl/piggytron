package handlers

import (
	"fmt"
	"net/http"

	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
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
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	banks, err := h.service.ReadAllBanksByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all banks", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var bviews []views.Bank
	for _, v := range banks {
		bviews = append(bviews, views.NewBank(v))
	}
	content := partials.Banks(bviews)
	renderWithMainLayout(w, r, "Banks", content)
}

func (h *BanksHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	aid := r.PathValue("id")
	bank, err := h.service.ReadOneByID(r.Context(), aid)
	if err != nil {
		logger.Error("error read one bank", "error", err, "account_id", aid)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<p>%s</p>", bank.Name())
	fmt.Fprintf(w, "<p>%s</p>", bank.Currency())
}
