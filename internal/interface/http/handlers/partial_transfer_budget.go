package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"

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
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
	"github.com/fdanctl/piggytron/web/views/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// TransferBudgetHandler transfers the budget from one category to another
type TransferBudgetHandler struct {
	service       *appbudget.Service
	categoryQuery query.CategoryQueryService
}

func NewTransferBudgetHandler(
	s *appbudget.Service,
	cq query.CategoryQueryService,
) *TransferBudgetHandler {
	return &TransferBudgetHandler{
		service:       s,
		categoryQuery: cq,
	}
}

func (h *TransferBudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get displays the transfer modal based transfer underspend or cover overspent
// type underspend || overspent
func (h *TransferBudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.URL.Query().Get("id")
	ttype := r.URL.Query().Get("type")
	month := r.URL.Query().Get("month")

	bm, err := budget.ParseMonth(month)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"TransferBudgetHandler.Get",
		)
		httperror.SendError(w, r, err)
		return
	}

	logger.Debug("type", "type", ttype)
	logger.Debug("month", "month", month)

	categoriesBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var opts []partials.BudgetSelectOption
	var totalAvailable int
	var catAvailable int
	for _, v := range categoriesBudgetSpent.Data {
		if v.Type != "income" {
			carryover := v.PrevTotalBudget - v.PrevTotalSpent
			available := v.Budgeted - v.Value + carryover
			totalAvailable += available
			if v.CategoryID != id {
				opts = append(opts, partials.BudgetSelectOption{
					Label:     v.Name,
					Value:     v.CategoryID,
					Type:      v.Type,
					Available: available,
				})
			} else {
				catAvailable = available
			}
		}
	}
	rta := categoriesBudgetSpent.Balance - totalAvailable - categoriesBudgetSpent.Hold
	opts = append([]partials.BudgetSelectOption{
		{
			Label:     "Ready to assign",
			Value:     "rta",
			Available: rta,
		},
	}, opts...)

	view := views.NewTransferBudgetForm()
	view.Amount = views.FormatAmount(math.Abs(float64(catAvailable)) / 100)
	view.Month = month
	var form templ.Component
	var title string
	if ttype == "underspend" {
		title = "Transfer underspend"
		view.FromID = id
		form = partials.TransferUnderspend(view, opts)
	} else {
		title = "Cover overspend"
		view.ToID = id
		form = partials.CoverOverspend(view, opts)
	}

	components.DialogWrapper(
		"",
		components.DialogHeader("", title, nil),
		form,
		nil,
		nil,
	).Render(r.Context(), w)
}

// Post transfers a budget from a one category to another
// ?from=<id>&to=<id>&amount=<number>
func (h *TransferBudgetHandler) Post(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	r.ParseForm()
	params := r.Form
	ttype := params.Get("type")
	amount := params.Get("amount")
	from := params.Get("from")
	to := params.Get("to")
	month := params.Get("month")
	punassign := params.Get("unassign")

	logger.Debug("form values", "amount", amount, "to", to, "from", from, "month", month)

	unassign, err := strconv.Atoi(punassign)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	excludeCatID := to
	if ttype == "underspend" {
		excludeCatID = from
	}

	bm, err := budget.ParseMonth(month)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"TransferBudgetHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	categoriesBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var selectOpts []partials.BudgetSelectOption
	var totalAvailable int
	for _, v := range categoriesBudgetSpent.Data {
		if v.Type != "income" {
			carryover := v.PrevTotalBudget - v.PrevTotalSpent
			available := v.Budgeted - v.Value + carryover
			totalAvailable += available
			if v.CategoryID != excludeCatID {
				selectOpts = append(selectOpts, partials.BudgetSelectOption{
					Label:     v.Name,
					Value:     v.CategoryID,
					Type:      v.Type,
					Available: available,
				})
			}
		}
	}
	rta := categoriesBudgetSpent.Balance - totalAvailable - categoriesBudgetSpent.Hold
	selectOpts = append([]partials.BudgetSelectOption{
		{
			Label:     "Ready to assign",
			Value:     "rta",
			Available: rta,
		},
	}, selectOpts...)

	view := views.TransferBudgetForm{
		Amount: amount,
		ToID:   to,
		FromID: from,
		Month:  month,
	}
	var form templ.Component
	if ttype == "underspend" {
		form = partials.TransferUnderspend(&view, selectOpts)
	} else {
		form = partials.CoverOverspend(&view, selectOpts)
	}
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
			"TransferBudgetHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	if err := h.service.TransferBudget(
		r.Context(),
		sessionInfo.UserID,
		from,
		to,
		bm,
		cents,
	); err != nil {
		view.SetError(err)
		if ttype == "underspend" {
			form = partials.TransferUnderspend(&view, selectOpts)
		} else {
			form = partials.CoverOverspend(&view, selectOpts)
		}
		httperror.SendFormError(w, r, err, form)
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

	budgett := unassign
	var budgeted int
	var changedCategoriesRows []templ.Component
	for i, v := range categoriesBudgetSpent.Data {
		var value int
		if v.Type == "income" {
			value = v.Value
			budgett += value
		} else {
			if v.CategoryID == to {
				categoriesBudgetSpent.Data[i].Budgeted += cents
				logger.Debug("to", "name", v.Name, "budgeted", v.Budgeted)
				carryover := v.PrevTotalBudget - v.PrevTotalSpent
				rowView := views.BudgetRowView{
					CategoryID: v.CategoryID,
					Month:      budget.Month(v.Month),
					Name:       v.Name,
					Carryover:  carryover,
					Budgeted:   categoriesBudgetSpent.Data[i].Budgeted,
					Available:  categoriesBudgetSpent.Data[i].Budgeted - v.Value + carryover,
				}
				changedCategoriesRows = append(
					changedCategoriesRows,
					pages.CategoryRow(v.Type, rowView, bm,
						templ.Attributes{
							"hx-swap-oob": "outerHTML",
						},
					),
				)
			}
			if v.CategoryID == from {
				categoriesBudgetSpent.Data[i].Budgeted -= cents
				logger.Debug("from", "name", v.Name, "budgeted", v.Budgeted)
				carryover := v.PrevTotalBudget - v.PrevTotalSpent
				rowView := views.BudgetRowView{
					CategoryID: v.CategoryID,
					Month:      budget.Month(v.Month),
					Name:       v.Name,
					Carryover:  carryover,
					Budgeted:   categoriesBudgetSpent.Data[i].Budgeted,
					Available:  categoriesBudgetSpent.Data[i].Budgeted - v.Value + carryover,
				}
				changedCategoriesRows = append(
					changedCategoriesRows,
					pages.CategoryRow(v.Type, rowView, bm,
						templ.Attributes{
							"hx-swap-oob": "outerHTML",
						},
					),
				)
			}
			value = categoriesBudgetSpent.Data[i].Budgeted
			budgeted += value
		}
		if value > 0 {
			node, link := charts.MakeBudgetSankeyNodeLink(v.Name, v.Type, value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	pageView := views.NewBudgetPageView(
		bm,
		categoriesBudgetSpent.MonthNet,
		categoriesBudgetSpent.Balance,
		categoriesBudgetSpent.Hold,
		categoriesBudgetSpent.Data,
	)

	ltb := budgett - budgeted - pageView.OnHold
	if ltb > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "Unassigned",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Budget",
				Target: "Unassigned",
				Value:  float32(ltb) / float32(100),
			},
		)
	}

	if pageView.OnHold > 0 {
		nodes = append(nodes, opts.SankeyNode{
			Name: "On hold",
			ItemStyle: &opts.ItemStyle{
				Color: "#D8DDF0",
			},
		})
		links = append(links,
			opts.SankeyLink{
				Source: "Budget",
				Target: "On hold",
				Value:  float32(pageView.OnHold) / float32(100),
			},
		)
	}

	theme := r.Header.Get("theme")
	component := components.NoData()
	if len(links) > 0 {
		sankey := charts.MakeSankey(nodes, links, true, theme)
		component = charts.ConvertChartToTemplComponent(sankey)
	}

	logger.Debug("overspent stat", "value", pageView.Overspent)

	w.Header().Set("HX-Trigger", "closeModal")

	joindedRows := templ.Join(
		changedCategoriesRows...,
	)
	obb := templ.Join(
		form,
		joindedRows,
		pages.BudgetStats(
			pageView.TotalBudgeted,
			pageView.ReadyToAssign,
			pageView.Income,
			pageView.UnassignCarryover,
			pageView.OnHold,
			pageView.AvailableToSpend,
			pageView.Overspent,
			pageView.Month,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.TotalRow("needs", pageView.NeedsBudget, pageView.NeedsAvailable,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.TotalRow("wants", pageView.WantsBudget, pageView.WantsAvailable,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.TotalRow(
			"savings",
			pageView.SavingsBudget,
			pageView.SavingsAvailable,
			templ.Attributes{
				"hx-swap-oob": "outerHTML",
			},
		),
		pages.PctSpan("needs", pageView.NeedsBudget, pageView.TotalBudgeted),
		pages.PctSpan("wants", pageView.WantsBudget, pageView.TotalBudgeted),
		pages.PctSpan("savings", pageView.SavingsBudget, pageView.TotalBudgeted),
		layouts.OOBWraper("budget-sankey", "innerHTML", nil, component),
	)

	obb.Render(r.Context(), w)
}
