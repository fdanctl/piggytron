package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/appledger"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// GoalContributeHandler renders the "contribute to goal" dialog (GET) and
// handles the contribution, which is a transfer into the goal account
// (POST).
type GoalContributeHandler struct {
	service       *appledger.Service
	categoryQuery query.CategoryQueryService
	accService    *appaccount.Service
}

func NewGoalContributeHandler(
	ts *appledger.Service,
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

// Get renders the contribution form pre-filled with the goal as
// destination (hidden) and category locked/disabled.
func (h *GoalContributeHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.NewTransferForm()
	acc, err := h.accService.FindOneByID(r.Context(), r.PathValue("id"), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	if acc.Type() != account.GoalType {
		err := errs.NewAppError(
			errs.KindNotFound,
			fmt.Sprintf("%s is not a goal", acc.ID()),
			fmt.Errorf("'%s' is not a goal: %w", acc.ID(), account.ErrAccountWrongType),
			"GoalContributeHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	_, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
		view.Date,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get categories select options: %w", err))
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
		view.Date,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to get account select options: %w", err))
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

	components.DialogWrapper("", components.DialogHeader("", "New Contribution", nil), form, nil, nil).
		Render(r.Context(), w)
}

// Post validates and creates the contribution transfer, triggering a
// transaction refetch and success toast.
func (h *GoalContributeHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	amount := r.FormValue("amount")
	currency := r.FormValue("currency")
	description := r.FormValue("description")
	date := r.FormValue("date")
	note := r.FormValue("note")
	category := r.FormValue("category")
	source := r.FormValue("source")
	destination := r.FormValue("destination")

	_, ecatOpts, err := getCategorySelectOptions(
		h.categoryQuery,
		r.Context(),
		sessionInfo.UserID,
		date,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
		date,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.TransferForm{
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
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amount),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amount, err),
			"GoalHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	d, err := time.Parse("02/01/2006", date)
	if err != nil {
		err := errs.NewGenericBadRequestAppError(err, "GoalHandler.Post")
		httperror.SendError(w, r, err)
		return
	}

	t, err := h.service.CreateTransfer(
		r.Context(),
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
		view.SetError(err)
		form = partials.GoalContributionForm(
			view,
			ecatOpts,
			append(noSavingsBanksOpts, goalSavingsOpts...),
		)
		httperror.SendFormError(w, r, err, form)
		return
	}

	logger.Debug(string(t.ID()))

	w.Header().Set("HX-Trigger", "refetch-transactions,closeModal")
	templ.Join(
		form,
		components.SendToast(components.Success, "Transfer transaction added"),
	).Render(r.Context(), w)
}
