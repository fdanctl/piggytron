package charts

import (
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeCatBarItems maps the per-month category values to a fixed 12-month
// slice of bar data (one entry per month, zero-filled for missing months).
func MakeCatBarItems(mvalues []query.CategoryMonthlyValue) []opts.BarData {
	byMonth := make(map[int]int, len(mvalues))
	for _, v := range mvalues {
		byMonth[v.Month] = v.Value
	}
	items := make([]opts.BarData, 0, 12)
	for m := 1; m <= 12; m++ {
		items = append(items, opts.BarData{Value: float64(byMonth[m]) / 100})
	}
	return items
}
