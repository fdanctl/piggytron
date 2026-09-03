package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// BankHandler renders the "new account" dialog (GET) and creates the bank
// account (POST).
type BankHandler struct {
	service          *appaccount.Service
	transactionQuery query.LedgerQueryService
}

func NewBankHandler(as *appaccount.Service, tq query.LedgerQueryService) *BankHandler {
	return &BankHandler{
		service:          as,
		transactionQuery: tq,
	}
}

func (h *BankHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")

	if action != "" {
		switch action {
		case "close":
			if r.Method != http.MethodPost {
				http.NotFound(w, r)
			}
			h.PostClose(w, r)

		case "initial-balance":
			switch r.Method {
			case http.MethodGet:
				h.GetChangeInitial(w, r)

			case http.MethodPost:
				h.PostChangeInitial(w, r)

			default:
				http.NotFound(w, r)
			}

		case "change-name":
			switch r.Method {
			case http.MethodGet:
				h.GetChangeName(w, r)

			case http.MethodPost:
				h.PostChangeName(w, r)

			default:
				http.NotFound(w, r)
			}
		}
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	case http.MethodDelete:
		h.Delete(w, r)

	default:
		http.NotFound(w, r)
	}
}

// Get renders the new-account form in a dialog.
func (h *BankHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := partials.BankForm(*views.NewBankForm())
	components.DialogWrapper(
		"",
		components.DialogHeader("", "New Account", nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

// Post validates and creates the bank, then navigates to the new bank page
// and prepends the account item to the accounts list.
func (h *BankHandler) Post(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.BankForm{
		Name:           r.FormValue("name"),
		Currency:       r.FormValue("currency"),
		Type:           r.FormValue("type"),
		InitialBalance: r.FormValue("initial-balance"),
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.BankForm(view).Render(r.Context(), w)
		return
	}

	cents, err := convertAmountStrToInt(view.InitialBalance)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", view.InitialBalance),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", view.InitialBalance, err),
			"BankHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	bank, err := h.service.CreateBank(
		r.Context(),
		sessionInfo.UserID,
		view.Name,
		view.Currency,
		view.Type,
		cents,
	)
	if err != nil {
		view.SetError(err)
		form := partials.BankForm(view)
		httperror.SendFormError(w, r, err, form)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/banks/%s",
			"transition": "true"
		}
		}`, bank.ID()),
	)

	bview := views.NewBank(
		string(bank.ID()),
		bank.Name(),
		string(bank.Type()),
		string(bank.Status()),
		0,
	)
	templ.Join(
		partials.BankForm(view),
		layouts.OOBWraper(
			"accounts-list",
			"beforeend",
			nil,
			partials.AccountItem(bview),
		),
	).Render(r.Context(), w)
}

// Delete deletes an account if no historical data is found.
func (h *BankHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		err := errs.NewAppError(
			errs.KindBadRequest,
			"Account ID is required",
			errors.New("no id passed"),
			"BankHandler.Delete",
		)
		httperror.SendError(w, r, err)
		return
	}

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	path := "/banks"
	if r.URL.Query().Get("type") == "goal" {
		path = "/goals"
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{"contentPush": { "url": "%s","transition": "true" }}`, path),
	)
}

// GetChangeName get the form to change name
func (h *BankHandler) GetChangeName(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	acc, err := h.service.FindOneByID(r.Context(), r.PathValue("id"), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.NewBankForm()
	view.Name = acc.Name()
	view.Currency = acc.Currency()
	view.Type = string(acc.Type())

	form := partials.BankNameForm(r.PathValue("id"), *view)
	components.DialogWrapper(
		"",
		components.DialogHeader("", "Change name", nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

// PostChangeName validates form changes an account name
func (h *BankHandler) PostChangeName(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.BankForm{
		Name:     r.FormValue("name"),
		Currency: r.FormValue("currency"),
		Type:     r.FormValue("type"),
	}

	err = h.service.UpdateBankName(
		r.Context(),
		sessionInfo.UserID,
		id,
		view.Name,
		view.Currency,
		view.Type,
	)
	if err != nil {
		view.SetError(err)
		form := partials.BankNameForm(id, view)
		httperror.SendFormError(w, r, err, form)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/banks/%s"
		}
		}`, id),
	)
	templ.Join(
		partials.BankNameForm(id, view),
		components.SendToast(components.Success, "Name changed with success"),
	).Render(r.Context(), w)
}

// PostClose closes an account name
func (h *BankHandler) PostClose(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.CloseAccount(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{"contentPush": { "url": "/banks/%s" }}`, id),
	)
}

func (h *BankHandler) GetChangeInitial(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")
	filters := query.NewLedgerFilters(
		[]string{"initial-balance"},
		[]string{id},
		nil,
		"",
		"",
		"",
		"",
	)

	var tid string
	var balance int
	initialBalanceEntry, err := h.transactionQuery.FindFiltered(
		r.Context(),
		sessionInfo.UserID,
		filters,
		1,
		0,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	if len(initialBalanceEntry) > 0 {
		balance = initialBalanceEntry[0].Amount
		tid = initialBalanceEntry[0].ID
	}

	view := views.NewInitialBalanceBankForm()
	view.InitialBalance = views.FormatAmount(float64(balance / 100))
	view.TransactionID = tid

	form := partials.BankInitialBalanceForm(id, *view)
	components.DialogWrapper(
		"",
		components.DialogHeader("", "Change initial balance", nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

func (h *BankHandler) PostChangeInitial(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.BankInitialBalanceForm{
		TransactionID:  r.FormValue("tid"),
		InitialBalance: r.FormValue("initial-balance"),
	}
	cents, err := convertAmountStrToInt(view.InitialBalance)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	err = h.service.UpdateBankInitial(
		r.Context(),
		sessionInfo.UserID,
		id,
		view.TransactionID,
		cents,
	)
	if err != nil {
		view.SetError(err)
		form := partials.BankInitialBalanceForm(id, view)
		httperror.SendFormError(w, r, err, form)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/banks/%s"
		}
		}`, id),
	)
	templ.Join(
		partials.BankInitialBalanceForm(id, view),
		components.SendToast(components.Success, "Initial balance changed with success"),
	).Render(r.Context(), w)
}
