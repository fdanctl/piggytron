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
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
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

	nodes := []opts.SankeyNode{
		{
			Name: "Budget",
			ItemStyle: &opts.ItemStyle{
				Color: "#194e4e",
			},
		},
	}
	var links []opts.SankeyLink
	if unassign > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "Unassigned Carryover",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Unassigned Carryover",
				Target: "Budget",
				Value:  float32(unassign) / float32(100),
			},
		)
	}

	for _, v := range categoryBudget {
		if v.Value > 0 {
			node, link := charts.MakeBudgetSankeyNodeLink(v.Name, v.Type, v.Value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	theme := r.Header.Get("theme")
	component := components.NoData()
	if len(links) > 0 {
		sankey := charts.MakeSankey(nodes, links, true, theme)
		component = charts.ConvertChartToTemplComponent(sankey)
	}

	obb := templ.Join(
		pages.BudgetInfoInputs(cents, month, cid),
		pages.CatRowLeftCell(cid, catLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
		pages.BudgetStats(
			totalBudgeted,
			leftToBudget,
			income,
			unassign,
			leftToSpent,
			overspent,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.TotalRow(catType, totalRowBudget, totalRowLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
		pages.PctSpan("needs", needsBudget, totalBudgeted),
		pages.PctSpan("wants", wantsBudget, totalBudgeted),
		pages.PctSpan("savings", savingsBudget, totalBudgeted),
		layouts.OOBWraper("budget-sankey", "innerHTML", nil, component),
	)

	obb.Render(r.Context(), w)
}
