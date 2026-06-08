package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/apptransaction"
	"github.com/fdanctl/piggytron/internal/domain/transaction"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalContributeHandler struct {
	service       *apptransaction.Service
	categoryQuery query.CategoryQueryService
	accService    *appaccount.Service
}

func NewGoalContributeHandler(
	ts *apptransaction.Service,
	cq query.CategoryQueryService,
	as *appaccount.Service,
) *GoalContributeHandler {
	return &GoalContributeHandler{
		service:       ts,
		categoryQuery: cq,
		accService:    as,
	}
}

func (h *GoalContributeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *GoalContributeHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	_, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		logger.Error("error reading all categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		logger.Error("error reading all accounts", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	view := views.NewTransferForm()
	acc, err := h.accService.FindOneByID(r.Context(), r.PathValue("id"))
	if err != nil {
		logger.Error("error finding goal ", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if acc.Type() != "goal" {
		logger.Error("error finding goal ", "error", err)
		http.Error(w, "internal error", http.StatusBadRequest)
		return
	}
	view.Description = fmt.Sprintf("%s contribution", acc.Name())
	view.DestinationAcc = r.PathValue("id")
	view.Category = string(*acc.CategoryID())
	form := partials.GoalContributionForm(
		*view,
		ecatOpts,
		append(noSavingsBanksOpts, goalSavingsOpts...),
	)

	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", "New Contribution", nil).Render(ctx, w)
}

func (h *GoalContributeHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	description := r.FormValue("description")
	date := r.FormValue("date")
	category := r.FormValue("category")
	source := r.FormValue("source")
	destination := r.FormValue("destination")

	_, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		logger.Error("error reading all categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		logger.Error("error reading all accounts", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

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

	form := partials.GoalContributionForm(
		view,
		ecatOpts,
		append(noSavingsBanksOpts, goalSavingsOpts...),
	)

	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)

		w.WriteHeader(http.StatusUnprocessableEntity)
		form.Render(r.Context(), w)
		return
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

	t, err := h.service.CreateTransfer(
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
		if errors.Is(err, transaction.ErrNegativeBalance) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			templ.Join(
				form,
				components.SendToast(components.Error, err.Error()),
			).Render(r.Context(), w)
			return
		}

		if errors.Is(err, transaction.ErrGoalCategory) {
			w.WriteHeader(http.StatusUnprocessableEntity)
			templ.Join(
				form,
				components.SendToast(components.Error, err.Error()),
			).Render(r.Context(), w)
			return
		}

		if errors.Is(err, transaction.ErrNotSavingsCategory) {
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

	logger.Debug(string(t.ID()))

	w.Header().Set("HX-Trigger", "refetch-transactions,closeModal")
	templ.Join(
		form,
		components.SendToast(components.Success, "Transfer transaction added"),
	).Render(r.Context(), w)
}
