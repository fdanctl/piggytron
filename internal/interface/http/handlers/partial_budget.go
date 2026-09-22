package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appbudget"
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// BudgetHandler handles budget amount edits on the budget page: it persists
// the new amount and re-renders every budget summary cell out-of-band.
type BudgetHandler struct {
	service       *appbudget.Service
	categoryQuery query.CategoryQueryService
}

func NewBudgetHandler(
	s *appbudget.Service,
	cq query.CategoryQueryService,
) *BudgetHandler {
	return &BudgetHandler{
		service:       s,
		categoryQuery: cq,
	}
}

func (h *BudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Post(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Post persists the new budget amount for the category and month and
// recomputes the whole budget summary (totals, percentages, sankey) from the
// form-provided deltas.
func (h *BudgetHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	r.ParseForm()
	params := r.Form
	amount := params.Get("amount")
	cid := params.Get("cid")
	month := params.Get("month")
	ps := params.Get("prev-amount")
	catType := params.Get("ctype")
	ptotalBudgeted := params.Get("total-budgeted")
	pneedstotalRowBudget := params.Get("needs-total-row-budget")
	pwantstotalRowBudget := params.Get("wants-total-row-budget")
	psavingstotalRowBudget := params.Get("savings-total-row-budget")
	ptotalRowLeft := params.Get("total-row-left")
	pcatLeft := params.Get("cat-left")
	pleftToBudget := params.Get("ltb")
	pleftToSpent := params.Get("lts")
	pincome := params.Get("income")
	ponHold := params.Get("on-hold")
	punassign := params.Get("unassign")
	poverspent := params.Get("overspent")

	prev, err := strconv.Atoi(ps)
	budgetInfoInputs := pages.BudgetInfoInputs(prev, month, cid)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	cents, err := convertAmountStrToInt(amount)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid amount", amount),
			fmt.Errorf("failed to convert amount '%s' to cents: %w", amount, err),
			"BudgetHandler.Post",
		)
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	bm := budget.NewMonth(time.Now())
	if month != "" {
		bm, err = budget.ParseMonth(month)
		if err != nil {
			err := errs.NewAppError(
				errs.KindBadRequest,
				fmt.Sprintf("%s is not a valid month", month),
				fmt.Errorf("failed to parse month '%s': %w", month, err),
				"BudgetHandler.Post",
			)
			httperror.SendError(w, r, err)
			return
		}
	}

	if cents == prev {
		logger.Debug("nothing to do")
		budgetInfoInputs.Render(r.Context(), w)
		return
	}

	_, err = h.service.SaveBudget(r.Context(), cid, bm, cents)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	addedAmount := cents - prev

	leftToSpent, err := strconv.Atoi(pleftToSpent)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	leftToSpent += addedAmount

	leftToBudget, err := strconv.Atoi(pleftToBudget)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	leftToBudget -= addedAmount

	income, err := strconv.Atoi(pincome)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	onHold, err := strconv.Atoi(ponHold)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	unassign, err := strconv.Atoi(punassign)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	totalBudgeted, err := strconv.Atoi(ptotalBudgeted)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	totalBudgeted += addedAmount

	overspent, err := strconv.Atoi(poverspent)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	catLeft, err := strconv.Atoi(pcatLeft)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	// if prev left is overspent reset this spent
	if catLeft < 0 {
		overspent += catLeft
	}
	// update
	catLeft += addedAmount
	// if new left is overspent add to overspent
	if catLeft < 0 {
		overspent -= catLeft
	}

	needsBudget, err := strconv.Atoi(pneedstotalRowBudget)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	wantsBudget, err := strconv.Atoi(pwantstotalRowBudget)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	savingsBudget, err := strconv.Atoi(psavingstotalRowBudget)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	var totalRowBudget int
	switch catType {
	case "needs":
		needsBudget += addedAmount
		totalRowBudget = needsBudget
	case "wants":
		wantsBudget += addedAmount
		totalRowBudget = wantsBudget
	case "savings":
		savingsBudget += addedAmount
		totalRowBudget = savingsBudget
	}

	totalRowLeft, err := strconv.Atoi(ptotalRowLeft)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	totalRowLeft += addedAmount

	categoryBudget, err := h.categoryQuery.GetCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		bm.Time(),
		time.Date(bm.Time().Year(), bm.Time().Month()+1, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		err := fmt.Errorf("error geting categories budget-spent: %w", err)
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}

	theme := r.Header.Get("theme")
	component := charts.BudgetChart(unassign, onHold, theme, categoryBudget)

	obb := templ.Join(
		pages.BudgetInfoInputs(cents, month, cid),
		layouts.HxPartial(
			fmt.Sprint("#ac-", cid),
			"outerHTML",
			pages.CatRowAvailableCell(cid, catLeft, bm, nil),
		),
		layouts.HxPartial(
			"#budget-stats",
			"outerHTML",
			pages.BudgetStats(
				totalBudgeted,
				leftToBudget,
				income,
				unassign,
				onHold,
				leftToSpent,
				overspent,
				bm,
				nil,
			),
		),
		layouts.HxPartial(
			fmt.Sprintf("#%s-total-row", catType),
			"outerHTML",
			pages.TotalRow(catType, totalRowBudget, totalRowLeft, nil),
		),
		pages.PctSpan("needs", needsBudget, totalBudgeted),
		pages.PctSpan("wants", wantsBudget, totalBudgeted),
		pages.PctSpan("savings", savingsBudget, totalBudgeted),
		layouts.HxPartial("#budget-sankey", "innerHTML", component),
	)

	obb.Render(r.Context(), w)
}
