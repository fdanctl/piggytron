package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type BankHandler struct {
	service *appaccount.Service
}

func NewBankHandler(as *appaccount.Service) *BankHandler {
	return &BankHandler{
		service: as,
	}
}

func (h *BankHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *BankHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := partials.BankForm(*views.NewBankForm())
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", "New Account", nil).Render(ctx, w)
}

func (h *BankHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	currency := r.FormValue("currency")
	savings := r.FormValue("savings")

	view := views.BankForm{
		Initial:   false,
		Name:      name,
		Currency:  currency,
		IsSavings: savings == "on",
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.BankForm(view).Render(r.Context(), w)
		return
	}

	bank, err := h.service.CreateBank(
		r.Context(),
		sessionInfo.UserID,
		name,
		currency,
		savings == "on",
	)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			logger.Info("invalid form - duplicated", "error", err)
			view.CustomError = err
		} else {
			logger.Error("error creating bank", "error", err)
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.BankForm(view).Render(r.Context(), w)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/banks/%s"
		}
		}`, bank.ID()),
	)

	templ.Join(
		partials.BankForm(view),
		layouts.OOBWraper(
			"accounts-list",
			"beforeend",
			nil,
			partials.AccountItem(string(bank.ID()), bank.Name()),
		),
	).Render(r.Context(), w)
}
