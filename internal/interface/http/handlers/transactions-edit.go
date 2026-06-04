package handlers

import (
	"errors"
	"net/http"
	"strconv"
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

type TransactionEditHandler struct {
	service      *transaction.Service
	inCatService *incomecategory.Service
	exCatService *expensecategory.Service
	accService   *accountapp.Service
}

func NewTransactionEditHandler(
	ts *transaction.Service,
	is *incomecategory.Service,
	es *expensecategory.Service,
	as *accountapp.Service,
) *TransactionEditHandler {
	return &TransactionEditHandler{
		service:      ts,
		inCatService: is,
		exCatService: es,
		accService:   as,
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
		logger.Error("unexpected error", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	id := r.PathValue("id")

	t, err := h.service.ReadOneByID(r.Context(), id)
	if err != nil {
		logger.Error("error finding transaction", "error", err)
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
		content = partials.IncomeForm(*v, icatOpts, accOptsNoGoals, string(t.ID()))

	case "expense":
		v := views.NewExpenseForm()
		v.Amount = strconv.Itoa(t.Amount()) // TODO fix
		v.Description = t.Description()
		v.Date = t.Date().Format("02/01/2006")
		v.Category = string(*t.ExpenseCategoryID())
		v.SourceAcc = string(*t.FromAccountID())
		title = "Edit Expense"
		content = partials.ExpenseForm(*v, ecatOpts, accOptsNoGoals, string(t.ID()))

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
		content = partials.TransferForm(*v, ecatOpts, accOpts, string(t.ID()))

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
		logger.Error("unexpected error", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
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

	var accOpts []components.SelectOption
	for _, v := range acc {
		accOpts = append(
			accOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
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
		form = partials.IncomeForm(view, icatOpts, accOpts, "")
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
		form = partials.ExpenseForm(view, ecatOpts, accOpts, "")
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

	// Update
	logger.Debug(id)
	logger.Debug(strconv.Itoa(cents))
	logger.Debug(d.Format(time.DateOnly))

	form.Render(r.Context(), w)
}

func (h *TransactionEditHandler) Delete(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())

	id := r.PathValue("id")

	err := h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, transactionDomain.ErrNegativeBalance) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			templ.Join(
				components.SendToast(components.Error, "Negative balance"),
			).Render(r.Context(), w)
			return
		}
		logger.Error("error creating transaction", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("HX-Trigger", "transaction-deleted")
}
