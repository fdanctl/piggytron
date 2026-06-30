package handlers

import (
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

type LedgerHandler struct {
	service       *appledger.Service
	categoryQuery query.CategoryQueryService
	accService    *appaccount.Service
}

func NewLedgerHandler(
	ts *appledger.Service,
	cq query.CategoryQueryService,
	as *appaccount.Service,
) *LedgerHandler {
	return &LedgerHandler{
		service:       ts,
		categoryQuery: cq,
		accService:    as,
	}
}

func (h *LedgerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *LedgerHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
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

	form := partials.TransactionForm(
		*views.NewIncomeForm(),
		*views.NewExpenseForm(),
		*views.NewTransferForm(),
		icatOpts,
		ecatOpts,
		append(noSavingsBanksOpts, goalSavingsOpts...),
		noSavingsBanksOpts,
	)

	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", "New Entry", nil).Render(ctx, w)
}

func (h *LedgerHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

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

	switch ttype {
	case "income":
		_, err := h.service.CreateIncome(
			r.Context(),
			sessionInfo.UserID,
			cents,
			currency,
			description,
			d,
			category,
			destination,
		)
		if err != nil {
			httperror.SendFormError(w, r, err, form)
			return
		}

	case "expense":
		_, err := h.service.CreateExpense(
			r.Context(),
			sessionInfo.UserID,
			cents,
			currency,
			description,
			d,
			category,
			source,
		)
		if err != nil {
			httperror.SendFormError(w, r, err, form)
			return
		}

	case "transfer":
		_, err := h.service.CreateTransfer(
			r.Context(),
			sessionInfo.UserID,
			cents,
			currency,
			description,
			d,
			category,
			source,
			destination,
		)
		if err != nil {
			httperror.SendFormError(w, r, err, form)
			return
		}
	}

	w.Header().Set("HX-Trigger", "refetch-transactions,closeModal")
	templ.Join(
		form,
		components.SendToast(
			components.Success,
			fmt.Sprintf("%s transaction added", views.CapitalizeFirst(ttype)),
		),
	).Render(r.Context(), w)
}
