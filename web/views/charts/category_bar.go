package charts

import (
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeCatBarItems makes per-month category values and month xAxis.
func MakeCatBarItems(mvalues []query.CategoryMonthlyValue) ([]opts.BarData, []string) {
	var items []opts.BarData
	var xAxis []string
	for _, v := range mvalues {
		xAxis = append(xAxis, budget.NewMonth(v.Month).Label())
		items = append(items, opts.BarData{Value: float64(v.Value) / 100})
	}
	return items, xAxis
}
