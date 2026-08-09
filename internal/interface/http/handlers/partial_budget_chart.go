package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/views/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BudgetChartHandler struct {
	categoryQuery query.CategoryQueryService
}

func NewBudgetChartHandler(
	cq query.CategoryQueryService,
) *BudgetChartHandler {
	return &BudgetChartHandler{
		categoryQuery: cq,
	}
}

func (h *BudgetChartHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BudgetChartHandler) Get(w http.ResponseWriter, r *http.Request) {
	month := r.PathValue("month")
	logger := middleware.LoggerFromContext(r.Context())
	logger.Debug(month)
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	punassign := r.URL.Query().Get("unassign")

	var unassignCarry int
	if punassign != "" {
		unassignCarry, err = strconv.Atoi(punassign)
		if err != nil {
			httperror.SendError(w, r, err)
			return
		}
	}

	bm := budget.NewMonth(time.Now())
	if month != "" {
		bm, err = budget.ParseMonth(month)
		if err != nil {
			err := errs.NewAppError(
				errs.KindBadRequest,
				fmt.Sprintf("%s is not a valid month", month),
				fmt.Errorf("failed to parse month '%s': %w", month, err),
				"BudgetChartHandler.Get",
			)
			httperror.SendError(w, r, err)
			return
		}
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

	nodes := []opts.SankeyNode{
		{
			Name: "Budget",
			ItemStyle: &opts.ItemStyle{
				Color: "#194e4e",
			},
		},
	}
	var links []opts.SankeyLink
	if unassignCarry > 0 {
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
				Value:  float32(unassignCarry) / float32(100),
			},
		)
	}

	budget := unassignCarry
	var budgeted int
	for _, v := range categoryBudget {
		if v.Type == "income" {
			budget += v.Value
		} else {
			budgeted += v.Value
		}
		if v.Value > 0 {
			node, link := charts.MakeBudgetSankeyNodeLink(v.Name, v.Type, v.Value)
			nodes = append(nodes, node)
			links = append(links, link)
		}
	}

	ltb := budget - budgeted
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

	theme := r.Header.Get("theme")

	component := components.NoData()
	if len(links) > 0 {
		sankey := charts.MakeSankey(nodes, links, true, theme)
		component = charts.ConvertChartToTemplComponent(sankey)
	}

	component.Render(r.Context(), w)
}
