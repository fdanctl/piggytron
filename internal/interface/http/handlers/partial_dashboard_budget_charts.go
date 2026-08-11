package handlers

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views/charts"
)

// DashboardBudgetCharts renders the dashboard's budget charts: the
// budget-vs-spent bar chart and the spent-by-type pie.
type DashboardBudgetCharts struct {
	categoryQuery query.CategoryQueryService
}

func NewDashboardBudgetCharts(
	cq query.CategoryQueryService,
) *DashboardBudgetCharts {
	return &DashboardBudgetCharts{
		categoryQuery: cq,
	}
}

func (h *DashboardBudgetCharts) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the budget bar chart and swaps the spent-by-type pie
// out-of-band for the current month.
// TODO: make it more complete,
//   - pie chart has a 2nd slide showing the spent progrees by type
//   - considering the pie chart to be a double donut with spent and budget
func (h *DashboardBudgetCharts) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	bm := budget.NewMonth(time.Now())
	categoryBudgetSpent, err := h.categoryQuery.GetCategoriesBudgetSpentValue(
		r.Context(),
		sessionInfo.UserID,
		bm,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	theme := r.Header.Get("theme")
	pieItems := charts.MakeSpentByTypePieItems(categoryBudgetSpent.Data)

	pie := components.NoData()
	if len(pieItems) > 0 {
		c := charts.PieRadius(pieItems, "Spent", theme)
		pie = charts.ConvertChartToTemplComponent(c)
	}

	budgetItems, spentItems, categories := charts.MakeBudgetSpentBarItems(
		categoryBudgetSpent.Data,
	)
	bar := charts.ConvertChartToTemplComponent(
		charts.CreateCategoryBudgetSpentBarChart(budgetItems, spentItems, categories, theme),
	)

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	nextMonthFirst := time.Date(today.Year(), today.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	monthLen := time.Date(
		nextMonthFirst.Year(),
		nextMonthFirst.Month(),
		nextMonthFirst.Day()-1,
		0,
		0,
		0,
		0,
		time.UTC,
	).Day()
	daysLeft := monthLen - today.Day() - 1

	spentCardSlides := []templ.Component{pages.ChartContainer(bar)}
	size := 4 // groups of 4
	for i := 0; i < len(categoryBudgetSpent.Data); i += size {
		end := i + size
		if end > len(categoryBudgetSpent.Data) {
			end = len(categoryBudgetSpent.Data)
		}

		var slide []templ.Component
		for k := i; k < end; k++ {
			c := categoryBudgetSpent.Data[k]
			if c.Type == "income" {
				i++
				end++
				continue
			}
			slide = append(
				slide,
				pages.CategorySpentProgress(c.Name, c.Type, c.Value, c.Budgeted, daysLeft),
			)
		}
		spentCardSlides = append(spentCardSlides, pages.CategorySpentContainer(slide))
	}

	templ.Join(
		components.CarouselCard(
			"",
			"Budget-spent by category",
			spentCardSlides,
			true,
			true,
			true,
			nil,
		),
		layouts.OOBWraper(
			"spent-by-type",
			"innerHTML",
			nil,
			pie,
		),
	).Render(r.Context(), w)
}
