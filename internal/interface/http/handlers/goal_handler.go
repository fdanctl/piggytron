package handlers

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalHandler struct {
	accService   *accountapp.Service
	exCatService *expensecategoryapp.Service
}

func NewGoalHandler(as *accountapp.Service, es *expensecategoryapp.Service) *GoalHandler {
	return &GoalHandler{
		accService:   as,
		exCatService: es,
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
	form := partials.GoalForm(*views.NewGoalForm())
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", nil).Render(ctx, w)
}

func (h *GoalHandler) Post(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := sessionInfoFormCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	currency := r.FormValue("currency")
	tamount := r.FormValue("target-amount")
	tdate := r.FormValue("target-date")
	cat := r.FormValue("category")

	view := views.GoalForm{
		Initial:      false,
		Name:         name,
		Currency:     currency,
		TargetAmount: tamount,
		TargetDate:   tdate,
		Category:     cat,
	}
	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalForm(view).Render(r.Context(), w)
		return
	}

	goal, err := h.accService.CreateGoal(
		r.Context(),
		sessionInfo.UserID,
		name,
		currency,
		tamount,
		tdate,
		cat,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.GoalForm(view).Render(r.Context(), w)
		return
	}
	fmt.Println(goal)
}
