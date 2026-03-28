package views

import (
	"fmt"
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

func formatMoney(amount float64, cur currency.Unit, lang language.Tag) string {
	p := message.NewPrinter(lang)
	symbol := currency.Symbol(cur)

	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	return p.Sprintf("%s%s%.2f", sign, symbol, amount)
}

func FormatFloat(x float64) string {
	s := fmt.Sprintf("%.2f", x)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func FormatDateMY(d time.Time) string {
	y, m, _ := d.Date()
	return fmt.Sprintf("%s %s", m, strconv.Itoa(y))
}

func FormatDate(d time.Time) string {
	return d.Format(time.DateOnly)
}
