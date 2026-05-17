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

func StringArrToStr(arr []string) string {
	if len(arr) > 0 {
		return arr[0]
	}
	return ""
}

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

func FormatAmount(v float64) string {
	p := message.NewPrinter(language.English)

	if math.Mod(v, 1) == 0 {
		// no decimals
		return p.Sprintf("%d", int64(v))
	}

	// keep 2 decimals
	return p.Sprintf("%.2f", v)
}

func FormatFloat(x float64) string {
	s := fmt.Sprintf("%.2f", x)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

// FormatDateMY returns for example June 2026
func FormatDateMY(d time.Time) string {
	y, m, _ := d.Date()
	return fmt.Sprintf("%s %s", m, strconv.Itoa(y))
}

// FormatDateOnly return for in YYYY-MM-DD format
func FormatDateOnly(d time.Time) string {
	return d.Format(time.DateOnly)
}

func FormatDate(date time.Time) string {
	y := date.Year()
	m := date.Month()
	d := date.Day()

	return fmt.Sprintf("%s %d, %d", m, d, y)
}

// capitalize first letter
func CapitalizeFirst(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(string(s[:1])) + s[1:]
}
