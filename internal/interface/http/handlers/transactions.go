package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	expensecategory "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategory "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/internal/application/transaction"
	transactionDomain "github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type TransactionHandler struct {
	service      *transaction.Service
	inCatService *incomecategory.Service
	exCatService *expensecategory.Service
	accService   *accountapp.Service
}

func NewTransactionHandler(
	ts *transaction.Service,
	is *incomecategory.Service,
	es *expensecategory.Service,
	as *accountapp.Service,
) *TransactionHandler {
	return &TransactionHandler{
		service:      ts,
		inCatService: is,
		exCatService: es,
		accService:   as,
	}
}

func (h *TransactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *TransactionHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	ecats, err := h.exCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var ecatOpts []components.SelectOption
	for _, v := range ecats {
		ecatOpts = append(
			ecatOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	icats, err := h.inCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all income categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	var icatOpts []components.SelectOption
	for _, v := range icats {
		icatOpts = append(
			icatOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	acc, err := h.accService.ReadAllByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all accounts", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var accOptsNoGoals []components.SelectOption
	var accOpts []components.SelectOption
	for _, v := range acc {
		accOpts = append(
			accOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
		if v.IsSaving() != nil && !*v.IsSaving() {
			accOptsNoGoals = append(
				accOptsNoGoals,
				components.SelectOption{Label: v.Name(), Value: string(v.ID())},
			)
		}
	}

	form := partials.TransactionForm(
		*views.NewIncomeForm(),
		*views.NewExpenseForm(),
		*views.NewTransferForm(),
		icatOpts,
		ecatOpts,
		accOpts,
		accOptsNoGoals,
	)

	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", nil).Render(ctx, w)
}

func (h *TransactionHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
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

	ecats, err := h.exCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var ecatOpts []components.SelectOption
	for _, v := range ecats {
		ecatOpts = append(
			ecatOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	icats, err := h.inCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	var icatOpts []components.SelectOption
	for _, v := range icats {
		icatOpts = append(
			icatOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	acc, err := h.accService.ReadAllByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all accounts", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var accOptsNoGoals []components.SelectOption
	var accOpts []components.SelectOption
	for _, v := range acc {
		accOpts = append(
			accOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
		if v.IsSaving() != nil && !*v.IsSaving() {
			accOptsNoGoals = append(
				accOptsNoGoals,
				components.SelectOption{Label: v.Name(), Value: string(v.ID())},
			)
		}
	}

	var form templ.Component
	switch ttype {
	case "income":
		view := views.IncomeForm{
			Initial:        false,
			Amount:         amount,
			Description:    description,
			Currency:       currency,
			Date:           date,
			Category:       category,
			DestinationAcc: destination,
		}
		form = partials.IncomeForm(view, icatOpts, accOptsNoGoals, "")
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	case "expense":
		view := views.ExpenseForm{
			Initial:     false,
			Amount:      amount,
			Description: description,
			Currency:    currency,
			Date:        date,
			Category:    category,
			SourceAcc:   source,
		}
		form = partials.ExpenseForm(view, ecatOpts, accOptsNoGoals, "")
		msgs := view.Validate()
		if len(msgs) > 0 {
			logger.Info("invalid form", "error", msgs)
			w.WriteHeader(http.StatusUnprocessableEntity)
			form.Render(r.Context(), w)
			return
		}

	case "transfer":
		view := views.TransferForm{
			Initial:        false,
			Amount:         amount,
			Description:    description,
			Currency:       currency,
			Date:           date,
			Category:       category,
			SourceAcc:      source,
			DestinationAcc: destination,
		}
		form = partials.TransferForm(view, ecatOpts, accOpts, "")
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
		http.Error(w, "invalid amount", http.StatusBadRequest)
		return
	}

	d, err := time.Parse("02/01/2006", date)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
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
			logger.Error("error creating transaction", "error", err)
			http.Error(w, "Internal error", http.StatusBadRequest)
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
			if errors.Is(err, transactionDomain.ErrNegativeBalance) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				templ.Join(
					form,
					components.SendToast(components.Error, "Negative balance"),
				).Render(r.Context(), w)
				return
			}
			logger.Error("error creating transaction", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
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
			if errors.Is(err, transactionDomain.ErrNegativeBalance) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				templ.Join(
					form,
					components.SendToast(components.Error, err.Error()),
				).Render(r.Context(), w)
				return
			}

			if errors.Is(err, transactionDomain.ErrGoalCategory) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				templ.Join(
					form,
					components.SendToast(components.Error, err.Error()),
				).Render(r.Context(), w)
				return
			}

			if errors.Is(err, transactionDomain.ErrNotSavingsCategory) {
				w.WriteHeader(http.StatusUnprocessableEntity)
				templ.Join(
					form,
					components.SendToast(
						components.Error,
						"Category must be savings type to send money to savings account",
					),
				).Render(r.Context(), w)
				return
			}

			logger.Error("error creating transaction", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
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
