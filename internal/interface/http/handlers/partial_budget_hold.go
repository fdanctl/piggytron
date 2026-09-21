package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appbudgethold"
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/domain/budgethold"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
	"github.com/fdanctl/piggytron/web/views/charts"
	"golang.org/x/text/currency"
	"golang.org/x/text/language"
)

// BudgetHoldHandler renders the "new income category" dialog (GET)
// and creates the category (POST).
type BudgetHoldHandler struct {
	service       *appbudgethold.Service
	categoryQuery query.CategoryQueryService
}

func NewBudgetHoldHandler(
	s *appbudgethold.Service,
	cq query.CategoryQueryService,
) *BudgetHoldHandler {
	return &BudgetHoldHandler{
		service:       s,
		categoryQuery: cq,
	}
}

func (h *BudgetHoldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the new income category form in a dialog.
func (h *BudgetHoldHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	month := r.URL.Query().Get("month")
	ltb := r.URL.Query().Get("ltb")
	logger.Debug("query values", "month", month, "ltb", ltb)

	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"BudgetHoldHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	leftToBudget, err := strconv.Atoi(ltb)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", ltb),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", ltb, err),
			"BudgetHoldHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	bm, err := budgethold.ParseMonth(month)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"BudgetHoldHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	view := views.NewBudgetHoldForm()

	hold, err := h.service.FindBudgetHold(r.Context(), sessionInfo.UserID, bm)
	if err != nil {
		if !errors.Is(err, budgethold.ErrNotFound) {
			httperror.SendError(w, r, err)
			return
		}
	}

	view.Month = bm.String()
	view.Amount = views.FormatAmount(float64(leftToBudget) / 100)
	if hold != nil && hold.Amount() != 0 {
		view.Amount = views.FormatAmount(float64(hold.Amount()) / 100)
	}

	form := partials.BudgetHoldForm(view)
	components.DialogWrapper(
		"",
		components.DialogHeader("", "Hold for next month", nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

// Post validates and creates the income category, appending the new item
// to the income categories list out-of-band.
func (h *BudgetHoldHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	r.ParseForm()
	params := r.Form
	month := params.Get("month")
	amount := params.Get("amount")
	view := views.BudgetHoldForm{
		Month:  month,
		Amount: amount,
	}
	logger.Debug("form values", "month", month, "amount", amount)

	msgs := view.Validate()
	if len(msgs) > 0 {
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.BudgetHoldForm(&view).Render(r.Context(), w)
		return
	}
	bm, err := budgethold.ParseMonth(month)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"BudgetHoldHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	cents, err := convertAmountStrToInt(amount)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amount),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amount, err),
			"BudgetHoldHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	hold, err := h.service.SaveBudgetHold(r.Context(), sessionInfo.UserID, bm, cents)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	// budget stats
	ptotalBudgeted := params.Get("total-budgeted")
	pleftToBudget := params.Get("ltb")
	pleftToSpent := params.Get("lts")
	pincome := params.Get("income")
	punassign := params.Get("unassign")
	poverspent := params.Get("overspent")
	ponHold := params.Get("on-hold")

	onHold, err := strconv.Atoi(ponHold)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	totalBudgeted, err := strconv.Atoi(ptotalBudgeted)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	leftToBudget, err := strconv.Atoi(pleftToBudget)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	leftToBudget -= cents - onHold

	income, err := strconv.Atoi(pincome)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	unassign, err := strconv.Atoi(punassign)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	leftToSpent, err := strconv.Atoi(pleftToSpent)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	overspent, err := strconv.Atoi(poverspent)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	categoryBudget, err := h.categoryQuery.GetCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		bm.Time(),
		time.Date(bm.Time().Year(), bm.Time().Month()+1, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		err := fmt.Errorf("error geting categories budget-spent: %w", err)
		httperror.SendError(w, r, err)
		return
	}

	theme := r.Header.Get("theme")
	component := charts.BudgetChart(unassign, hold.Amount(), theme, categoryBudget)

	w.Header().Set("HX-Trigger", "closeModal")

	templ.Join(
		partials.BudgetHoldForm(&view),
		pages.BudgetStats(
			totalBudgeted,
			leftToBudget,
			income,
			unassign,
			hold.Amount(),
			leftToSpent,
			overspent,
			budget.NewMonth(bm.Time()),
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		components.SendToast(
			components.Success,
			fmt.Sprintf(
				"%s on hold for next month",
				views.FormatMoney(
					float64(hold.Amount())/100,
					currency.EUR,
					language.AmericanEnglish,
				),
			),
		),
		layouts.OOBWraper("budget-sankey", "innerHTML", nil, component),
	).Render(r.Context(), w)
}
