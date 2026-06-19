package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/application/apptransaction"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
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
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	_, ecatOpts, err := getCategorySelectOptions(
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
		httperror.SendError(w, r, err)
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
		httperror.SendError(w, r, err)
		return
	}

	noSavingsBanksOpts, goalSavingsOpts, err := getAccSelectOptions(
		h.accService,
		r.Context(),
		sessionInfo.UserID,
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
