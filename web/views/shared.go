// Package views builds the view models rendered by the templ templates: page
// data (e.g. BankPage) and the money/date formatting helpers shared across
// pages.
package views

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/currency"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// StringArrToStr returns the first element of a string slice, or "" when
// empty; used to read the first value of multi-value query strings.
func StringArrToStr(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

// FormatMoney renders an amount with the currency symbol and two decimals,
// localized for the given language tag.
func FormatMoney(amount float64, cur currency.Unit, lang language.Tag) string {
	p := message.NewPrinter(lang)
	symbol := currency.Symbol(cur)

	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return p.Sprintf("%s%s%.2f", sign, symbol, amount)
}

// FormatAmount renders a float without decimals when whole, otherwise with
// two decimals.
func FormatAmount(v float64) string {
	p := message.NewPrinter(language.English)

	if math.Mod(v, 1) == 0 {
		// no decimals
		return p.Sprintf("%d", int64(v))
	}

	// keep 2 decimals
	return p.Sprintf("%.2f", v)
}

// FormatFloat renders a float with up to two decimals, trimming trailing
// zeros.
func FormatFloat(x float64) string {
	s := fmt.Sprintf("%.2f", x)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// FormatDateMY renders a date as "Month YYYY" (eg. "June 2026").
func FormatDateMY(d time.Time) string {
	y, m, _ := d.Date()
	return fmt.Sprintf("%s %s", m, strconv.Itoa(y))
}

// FormatDateOnly return for in YYYY-MM-DD format
func FormatDateOnly(d time.Time) string {
	return d.Format(time.DateOnly)
}

// FormatDate renders a date as "Month D, YYYY" (e.g. "June 26, 2026").
func FormatDate(date time.Time) string {
	y := date.Year()
	m := date.Month()
	d := date.Day()

	return fmt.Sprintf("%s %d, %d", m, d, y)
}

// CapitalizeFirst capitalizes the first letter of a string
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[:1])) + s[1:]
}

// MonthDiff returns the number of whole calendar months between a and b.
func MonthDiff(a, b time.Time) int {
	return (b.Year()-a.Year())*12 + int(b.Month()-a.Month())
}
