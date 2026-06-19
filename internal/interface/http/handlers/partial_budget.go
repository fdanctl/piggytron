package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appbudget"
	"github.com/fdanctl/piggytron/internal/application/appcharts"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/internal/util"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BudgetHandler struct {
	service       *appbudget.Service
	chartsService *appcharts.Service
	categoryQuery query.CategoryQueryService
}

func NewBudgetHandler(
	s *appbudget.Service,
	cs *appcharts.Service,
	cq query.CategoryQueryService,
) *BudgetHandler {
	return &BudgetHandler{
		service:       s,
		chartsService: cs,
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
	bid := params.Get("bid")
	ps := params.Get("prev-amount")
	catType := params.Get("ctype")
	ptotalBudgeted := params.Get("total-budgeted")
	ptotalRowBudget := params.Get("total-row-budget")
	ptotalRowLeft := params.Get("total-row-left")
	pcatLeft := params.Get("cat-left")
	pleftToBudget := params.Get("ltb")
	pleftToSpent := params.Get("lts")
	pincome := params.Get("income")
	poverspent := params.Get("overspent")

	budgetID := bid
	prev, err := strconv.Atoi(ps)
	budgetInfoInputs := pages.BudgetInfoInputs(prev, budgetID, cid)
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

	now := time.Now()

	if cents == prev {
		logger.Debug("nothing to do")
		budgetInfoInputs.Render(r.Context(), w)
		return
	}

	if bid == "" || bid == util.ZeroUUID {
		_, err := h.service.CreateBudget(r.Context(), sessionInfo.UserID, cid, now, cents)
		if err != nil {
			httperror.SendFormError(w, r, err, budgetInfoInputs)
			return
		}
	} else {
		err := h.service.UpdateBudgetAmount(r.Context(), bid, cents)
		if err != nil {
			httperror.SendFormError(w, r, err, budgetInfoInputs)
			return
		}
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

	totalRowBudget, err := strconv.Atoi(ptotalRowBudget)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	totalRowBudget += addedAmount

	totalRowLeft, err := strconv.Atoi(ptotalRowLeft)
	if err != nil {
		httperror.SendFormError(w, r, err, budgetInfoInputs)
		return
	}
	totalRowLeft += addedAmount

	minD := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	maxD := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	categoryBudget, err := h.categoryQuery.GetCategoriesBudgetSpent(
		r.Context(),
		sessionInfo.UserID,
		minD,
		maxD,
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

	for _, v := range categoryBudget {
		if v.Value > 0 {
			node, link := h.chartsService.MakeBudgetSankeyNodeLink(v.Name, v.Type, v.Value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	sankey := h.chartsService.MakeSankey(nodes, links, false)
	chartComponent := h.chartsService.ConvertChartToTemplComponent(sankey)

	obb := templ.Join(
		pages.BudgetInfoInputs(cents, budgetID, cid),
		pages.CatRowLeftCell(cid, catLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
		pages.BudgetStats(
			totalBudgeted,
			leftToBudget,
			income,
			leftToSpent,
			overspent,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.TotalRow(catType, totalRowBudget, totalRowLeft, templ.Attributes{
			"hx-swap-oob": "outerHTML",
		}),
		pages.PctSpan(catType, totalRowBudget, totalBudgeted),
		layouts.OOBWraper("budget-sankey", "innerHTML", nil, chartComponent),
	)

	obb.Render(r.Context(), w)
}
