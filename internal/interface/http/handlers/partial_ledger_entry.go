package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/appledger"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// LedgerEntryHandler renders the edit dialog for one entry (GET) and
// handles its update (PUT) and deletion (DELETE).
type LedgerEntryHandler struct {
	service       *appledger.Service
	categoryQuery query.CategoryQueryService
	accService    *appaccount.Service
}

func NewLedgerEntryHandler(
	ts *appledger.Service,
	cq query.CategoryQueryService,
	as *appaccount.Service,
) *LedgerEntryHandler {
	return &LedgerEntryHandler{
		service:       ts,
		categoryQuery: cq,
		accService:    as,
	}
}

func (h *LedgerEntryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPut:
		h.Put(w, r)

	case http.MethodDelete:
		h.Delete(w, r)

	default:
		http.NotFound(w, r)
	}
}

// Get renders the edit dialog pre-filled from the entry's current values.
func (h *LedgerEntryHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")

	t, err := h.service.FindOneByID(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	icatOpts, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
		"",
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get categories select options: %w", err))
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
		"",
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get account select options: %w", err))
		return
	}

	var title string
	var content templ.Component
	var note string
	if t.Note() != nil {
		note = *t.Note()
	}
	switch string(t.Type()) {
	case "income":
		v := views.NewIncomeForm()
		v.Amount = views.FormatAmount(float64(t.Amount() / 100))
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.IncomeCategoryID())
		v.DestinationAcc = string(*t.ToAccountID())
		v.Note = note
		title = "Edit Income"
		content = partials.IncomeForm(*v, icatOpts, noSavingsBanksOpts, string(t.ID()))

	case "expense":
		v := views.NewExpenseForm()
		v.Amount = views.FormatAmount(float64(t.Amount() / 100))
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.ExpenseCategoryID())
		v.SourceAcc = string(*t.FromAccountID())
		v.Note = note
		title = "Edit Expense"
		content = partials.ExpenseForm(*v, ecatOpts, noSavingsBanksOpts, string(t.ID()))

	case "transfer":
		v := views.NewTransferForm()
		v.Amount = views.FormatAmount(float64(t.Amount() / 100))
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		if t.ExpenseCategoryID() != nil {
			v.Category = string(*t.ExpenseCategoryID())
		}
		v.DestinationAcc = string(*t.ToAccountID())
		v.SourceAcc = string(*t.FromAccountID())
		v.Note = note
		title = "Edit Transfer"
		content = partials.TransferForm(
			*v,
			ecatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			string(t.ID()),
		)

	case "interest":
		v := views.NewInterestForm()
		v.Amount = views.FormatAmount(float64(t.Amount() / 100))
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.IncomeCategoryID())
		v.DestinationAcc = string(*t.ToAccountID())
		v.Note = note
		title = "Edit Interest"
		content = partials.InterestForm(
			*v,
			icatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			string(t.ID()),
		)

	default:
		logger.Debug("DEFAULT")
	}

	components.DialogWrapper(
		"",
		components.DialogHeader("", title, nil),
		content,
		partials.FormActionBtns("Confirm"),
		nil,
	).Render(r.Context(), w)
}

// Put validates and applies the updated entry via the ledger service,
// re-rendering the form on error or showing a success toast.
func (h *LedgerEntryHandler) Put(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")
	ttype := r.FormValue("type")
	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	description := r.FormValue("description")
	date := r.FormValue("date")
	note := r.FormValue("note")
	category := r.FormValue("category")
	source := r.FormValue("source")
	destination := r.FormValue("destination")

	icatOpts, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
		date,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get categories select options: %w", err))
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
		date,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get account select options: %w", err))
		return
	}

	var form templ.Component
	switch ttype {
	case "income":
		view := views.IncomeForm{
			Amount:         amount,
			Description:    description,
			Currency:       currency,
			Date:           date,
			Category:       category,
			DestinationAcc: destination,
			Note:           note,
		}
		form = partials.IncomeForm(view, icatOpts, noSavingsBanksOpts, id)
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	case "expense":
		view := views.ExpenseForm{
			Amount:      amount,
			Description: description,
			Currency:    currency,
			Date:        date,
			Category:    category,
			SourceAcc:   source,
			Note:        note,
		}
		form = partials.ExpenseForm(view, ecatOpts, noSavingsBanksOpts, id)
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	case "transfer":
		view := views.TransferForm{
			Amount:         amount,
			Description:    description,
			Currency:       currency,
			Date:           date,
			Category:       category,
			SourceAcc:      source,
			DestinationAcc: destination,
			Note:           note,
		}
		form = partials.TransferForm(
			view,
			ecatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			id,
		)
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	case "interest":
		view := views.InterestForm{
			Amount:         amount,
			Description:    description,
			Currency:       currency,
			Date:           date,
			Category:       category,
			DestinationAcc: destination,
			Note:           note,
		}
		form = partials.InterestForm(
			view,
			icatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			id,
		)
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	default:
		httperror.SendError(
			w,
			r,
			errs.NewGenericBadRequestAppError(errors.New("invalid type"), "LedgerEntryHandler.Put"),
		)
		return
	}

	cents, err := convertAmountStrToInt(amount)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amount),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amount, err),
			"LedgerEntryHandler.Put",
		)
		httperror.SendError(w, r, err)
		return
	}

	d, err := time.Parse("02/01/2006", date)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid date", date),
			fmt.Errorf("failed to parse date '%s': %w", date, err),
			"LedgerEntryHandler.Put",
		)
		httperror.SendError(w, r, err)
		return
	}

	_, err = h.service.Update(
		r.Context(),
		ttype,
		id,
		sessionInfo.UserID,
		cents,
		currency,
		description,
		d,
		note,
		category,
		source,
		destination,
	)
	if err != nil {
		httperror.SendFormError(w, r, err, form)
		return
	}

	w.Header().Set("HX-Trigger", "refetch-transactions,closeAllModal")
	templ.Join(
		form,
		components.SendToast(
			components.Success,
			fmt.Sprintf("%s transaction updated", views.CapitalizeFirst(ttype)),
		),
	).Render(r.Context(), w)
}

// Delete removes the entry and triggers the transaction-deleted event.
func (h *LedgerEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	w.Header().Set("HX-Trigger", "transaction-deleted")
	components.SendToast(
		components.Success,
		"Transaction deleted",
	).Render(r.Context(), w)
}
