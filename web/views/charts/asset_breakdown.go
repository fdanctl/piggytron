package charts

import (
	"slices"

	"github.com/fdanctl/piggytron/internal/query"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// MakeAccountsPieItems converts account balances into pie chart data,
// sorted from largest to smallest balance.
func MakeAccountsPieItems(acc []query.AccountWithSum) []opts.PieData {
	slices.SortFunc(acc, func(a, b query.AccountWithSum) int {
		return b.Sum - a.Sum
	})
	var data []opts.PieData
	for _, v := range acc {
		data = append(data, opts.PieData{Name: v.Name, Value: float64(v.Sum) / float64(100)})
	}
	return data
}
