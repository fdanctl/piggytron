package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalHandler struct {
	accService          *appaccount.Service
	accountQueryService query.AccountQueryService
	categoryQuery       query.CategoryQueryService
}

func NewGoalHandler(
	as *appaccount.Service,
	aq query.AccountQueryService,
	cq query.CategoryQueryService,
) *GoalHandler {
	return &GoalHandler{
		accService:          as,
		accountQueryService: aq,
		categoryQuery:       cq,
	}
}

func (h *GoalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	case http.MethodPut:
		h.Put(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *GoalHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	view := views.NewGoalForm()
	id := r.PathValue("id")
	title := "New Goal"
	if id != "" {
		g, err := h.accService.FindOneByID(r.Context(), id, sessionInfo.UserID)
		if err != nil {
			logger.Error("error finding goal", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		if g.Type() != account.GoalType {
			http.Error(w, "Not a goal", http.StatusUnprocessableEntity)
			return
		}
		view.Name = g.Name()
		view.Currency = g.Currency()
		view.StartDate = g.StartDate().Format("02/01/2006")
		view.TargetAmount = views.FormatAmount((float64(*g.TargetAmount()) / float64(100)))
		view.Category = string(*g.CategoryID())
		if g.TargetDate() != nil {
			view.TargetDate = g.TargetDate().Format("02/01/2006")
		}
		title = "Edit Goal"
	}

	form := partials.GoalForm(*view, ecatOpts, id)
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", title, nil).Render(ctx, w)
}

func (h *GoalHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	formData := h.parseGoalForm(r)

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

	view := views.GoalForm{
		Initial:      false,
		Name:         formData.name,
		Currency:     formData.currency,
		TargetAmount: formData.tamount,
		StartDate:    formData.sdate,
		TargetDate:   formData.tdate,
		Category:     formData.cat,
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalFormContent(view, ecatOpts).Render(r.Context(), w)
		return
	}
	amount, startDate, pDate, err := h.parseGoalFormValues(view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goal, err := h.accService.CreateGoal(
		r.Context(),
		sessionInfo.UserID,
		formData.name,
		formData.currency,
		amount,
		startDate,
		pDate,
		formData.cat,
	)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			logger.Info("invalid form - duplicated", "error", err)
			view.CustomError = err
		} else {
			logger.Error("error creating goal", "error", err)
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalForm(view, ecatOpts, "").Render(r.Context(), w)
		return
	}

	g, err := h.accountQueryService.FindWithSum(
		r.Context(), string(goal.ID()),
	)
	if err != nil {
		logger.Error("error finding accounts", "error", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/goals/%s",
			"transition": "true"
		}
		}`, goal.ID()),
	)

	templ.Join(
		partials.GoalFormContent(view, ecatOpts),
		layouts.OOBWraper(
			"active-goals-list",
			"beforeend",
			nil,
			partials.GoalItem(views.NewGoal(*g), nil),
		),
	).Render(r.Context(), w)
}

func (h *GoalHandler) Put(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "Goal ID is required", http.StatusBadRequest)
		return
	}
	formData := h.parseGoalForm(r)

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

	view := views.GoalForm{
		Initial:      false,
		Name:         formData.name,
		Currency:     formData.currency,
		TargetAmount: formData.tamount,
		StartDate:    formData.sdate,
		TargetDate:   formData.tdate,
		Category:     formData.cat,
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalFormContent(view, ecatOpts).Render(r.Context(), w)
		return
	}

	amount, startDate, pDate, err := h.parseGoalFormValues(view)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goal, err := h.accService.UpdateGoal(
		r.Context(),
		id,
		sessionInfo.UserID,
		formData.name,
		formData.currency,
		amount,
		startDate,
		pDate,
		formData.cat,
	)
	if err != nil {
		if errors.Is(err, account.ErrContributionBeforeStartDate) {
			view.CustomError = err
		} else {
			logger.Error("error updating goal", "error", err)
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalFormContent(view, ecatOpts).Render(r.Context(), w)
		return
	}

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeAllModal": true,
		"contentPush": {
			"url": "/goals/%s"
		}
		}`, goal.ID()),
	)

	partials.GoalFormContent(view, ecatOpts).Render(r.Context(), w)
}

type goalFormData struct {
	name, currency, tamount, sdate, tdate, cat string
}

func (h *GoalHandler) parseGoalForm(r *http.Request) goalFormData {
	return goalFormData{
		name:     r.FormValue("name"),
		currency: r.FormValue("currency"),
		tamount:  r.FormValue("target-amount"),
		sdate:    r.FormValue("start-date"),
		tdate:    r.FormValue("target-date"),
		cat:      r.FormValue("category"),
	}
}

func (h *GoalHandler) parseGoalFormValues(
	view views.GoalForm,
) (amount int, startDate time.Time, pTargetDate *time.Time, err error) {
	amount, err = convertAmountStrToInt(view.TargetAmount)
	if err != nil {
		err = errors.New("Invalid amount")
		return
	}

	startDate, err = time.Parse("02/01/2006", view.StartDate)
	if err != nil {
		err = errors.New("Invalid date")
		return
	}

	targetDate, err := time.Parse("02/01/2006", view.TargetDate)
	if err == nil {
		pTargetDate = &targetDate
	}
	err = nil
	return
}
