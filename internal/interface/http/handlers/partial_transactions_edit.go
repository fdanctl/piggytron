package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/apptransaction"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type TransactionEditHandler struct {
	service       *apptransaction.Service
	categoryQuery query.CategoryQueryService
	accService    *appaccount.Service
}

func NewTransactionEditHandler(
	ts *apptransaction.Service,
	cq query.CategoryQueryService,
	as *appaccount.Service,
) *TransactionEditHandler {
	return &TransactionEditHandler{
		service:       ts,
		categoryQuery: cq,
		accService:    as,
	}
}

func (h *TransactionEditHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

func (h *TransactionEditHandler) Get(w http.ResponseWriter, r *http.Request) {
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
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get categories select options: %w", err))
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get account select options: %w", err))
		return
	}

	var title string
	var content templ.Component
	switch string(t.Type()) {
	case "income":
		v := views.NewIncomeForm()
		v.Amount = strconv.Itoa(t.Amount()) // TODO fix
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.IncomeCategoryID())
		v.DestinationAcc = string(*t.ToAccountID())
		title = "Edit Income"
		content = partials.IncomeForm(*v, icatOpts, noSavingsBanksOpts, string(t.ID()))

	case "expense":
		v := views.NewExpenseForm()
		v.Amount = strconv.Itoa(t.Amount()) // TODO fix
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.ExpenseCategoryID())
		v.SourceAcc = string(*t.FromAccountID())
		title = "Edit Expense"
		content = partials.ExpenseForm(*v, ecatOpts, noSavingsBanksOpts, string(t.ID()))

	case "transfer":
		v := views.NewTransferForm()
		v.Amount = strconv.Itoa(t.Amount()) // TODO fix
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		if t.ExpenseCategoryID() != nil {
			v.Category = string(*t.ExpenseCategoryID())
		}
		v.DestinationAcc = string(*t.ToAccountID())
		v.SourceAcc = string(*t.FromAccountID())
		title = "Edit Transfer"
		content = partials.TransferForm(
			*v,
			ecatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			string(t.ID()),
		)

	default:
		logger.Debug("DEFAULT")
	}

	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("", title, nil).Render(ctx, w)
}

func (h *TransactionEditHandler) Put(w http.ResponseWriter, r *http.Request) {
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
	category := r.FormValue("category")
	source := r.FormValue("source")
	destination := r.FormValue("destination")

	icatOpts, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get categories select options: %w", err))
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
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
		}
		form = partials.IncomeForm(view, icatOpts, noSavingsBanksOpts, "")
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
		}
		form = partials.ExpenseForm(view, ecatOpts, noSavingsBanksOpts, "")
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
		}
		form = partials.TransferForm(
			view,
			ecatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
			"",
		)
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	default:
		logger.Debug("DEFAULT")
	}

	cents, err := convertAmountStrToInt(amount)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amount),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amount, err),
			"BudgetHandler.Post",
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
			"BudgetHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	// Update
	logger.Debug(id)
	logger.Debug(strconv.Itoa(cents))
	logger.Debug(d.Format(time.DateOnly))

	form.Render(r.Context(), w)
}

func (h *TransactionEditHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	w.Header().Set("HX-Trigger", "transaction-deleted")
}
