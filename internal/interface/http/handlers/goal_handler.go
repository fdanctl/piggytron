package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalHandler struct {
	accService          *accountapp.Service
	accountQueryService query.AccountQueryService
	exCatService        *expensecategoryapp.Service
}

func NewGoalHandler(
	as *accountapp.Service,
	aq query.AccountQueryService,
	es *expensecategoryapp.Service,
) *GoalHandler {
	return &GoalHandler{
		accService:          as,
		accountQueryService: aq,
		exCatService:        es,
	}
}

func (h *GoalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

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

	cats, err := h.exCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var catOpts []components.SelectOption
	for _, v := range cats {
		catOpts = append(
			catOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	form := partials.GoalForm(*views.NewGoalForm(), catOpts)
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", nil).Render(ctx, w)
}

func (h *GoalHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	currency := r.FormValue("currency")
	tamount := r.FormValue("target-amount")
	sdate := r.FormValue("start-date")
	tdate := r.FormValue("target-date")
	cat := r.FormValue("category")

	cats, err := h.exCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all expense categories", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var catOpts []components.SelectOption
	for _, v := range cats {
		catOpts = append(
			catOpts,
			components.SelectOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	view := views.GoalForm{
		Initial:      false,
		Name:         name,
		Currency:     currency,
		TargetAmount: tamount,
		StartDate:    sdate,
		TargetDate:   tdate,
		Category:     cat,
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalForm(view, catOpts).Render(r.Context(), w)
		return
	}

	amount, err := convertAmountStrToInt(tamount)
	if err != nil {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("02/01/2006", sdate)
	if err != nil {
		http.Error(w, "Invalid date", http.StatusBadRequest)
		return
	}

	targetDate, err := time.Parse("02/01/2006", tdate)
	var pDate *time.Time
	if err == nil {
		pDate = &targetDate
	}

	goal, err := h.accService.CreateGoal(
		r.Context(),
		sessionInfo.UserID,
		name,
		currency,
		amount,
		startDate,
		pDate,
		cat,
	)
	if err != nil {
		if errors.Is(err, account.ErrDuplicate) {
			logger.Info("invalid form - duplicated", "error", err)
			view.CustomError = err
		} else {
			logger.Error("error creating goal", "error", err)
		}

		cats, err := h.exCatService.ReadAllUserCategories(r.Context(), sessionInfo.UserID)
		if err != nil {
			logger.Error("error reading all expense categories", "error", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		var catOpts []components.SelectOption
		for _, v := range cats {
			catOpts = append(
				catOpts,
				components.SelectOption{Label: v.Name(), Value: string(v.ID())},
			)
		}

		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalForm(view, catOpts).Render(r.Context(), w)
		return
	}

	g, err := h.accountQueryService.FindWithSum(
		r.Context(), string(goal.ID()),
	)

	w.Header().Set(
		"HX-Trigger",
		fmt.Sprintf(`{
		"closeModal": true,
		"contentPush": {
			"url": "/goals/%s"
		}
		}`, goal.ID()),
	)

	templ.Join(
		partials.GoalForm(view, catOpts),
		layouts.OOBWraper(
			"active-goals-list",
			"beforeend",
			nil,
			partials.GoalItem(views.NewGoal(*g), nil),
		),
	).Render(r.Context(), w)
}
