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
	"unicode/utf8"

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

// ConvertAmountStrToInt converts user input to amount in cents.
func ConvertAmountStrToInt(str string) (int, error) {
	str = strings.ReplaceAll(str, ",", "")
	i := strings.Index(str, ".")
	tAmount := 0

	length := utf8.RuneCountInString(str)
	if str == "" {
		return 0, nil
	}

	if i == -1 {
		parsed, err := strconv.Atoi(str)
		if err != nil {
			return 0, err
		}
		return parsed * 100, nil
	}

	if length-1-i > 2 {
		return 0, nil
	}

	for length-i < 3 {
		str += "0"
		length++
	}

	tAmount, err := strconv.Atoi(strings.Replace(str, ".", "", 1))
	if err != nil {
		return 0, err
	}

	return tAmount, nil
}

// DateDiffString calculates the diference between two dates and
// returns a string of "%d years, %d months, %d days"
func DateDiffString(start, end time.Time) string {
	if end.Before(start) {
		start, end = end, start
	}

	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())
	days := end.Day() - start.Day()

	if days < 0 {
		months--

		// Number of days in the month before end.
		prevMonth := end.AddDate(0, 0, -end.Day())
		days += prevMonth.Day()
	}

	if months < 0 {
		years--
		months += 12
	}

	var strs []string

	if years > 0 {
		if years == 1 {
			strs = append(strs, fmt.Sprintf("%d year", years))
		} else {
			strs = append(strs, fmt.Sprintf("%d years", years))
		}
	}

	if months > 0 {
		if months == 1 {
			strs = append(strs, fmt.Sprintf("%d month", months))
		} else {
			strs = append(strs, fmt.Sprintf("%d months", months))
		}
	}

	if days > 0 {
		if days == 1 {
			strs = append(strs, fmt.Sprintf("%d day", days))
		} else {
			strs = append(strs, fmt.Sprintf("%d days", days))
		}
	}

	return strings.Join(strs, ", ")
}
