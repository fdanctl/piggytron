package budget

import (
	"fmt"
	"time"
)

var monthAbbr = map[time.Month]string{
	time.January:   "Jan",
	time.February:  "Feb",
	time.March:     "Mar",
	time.April:     "Apr",
	time.May:       "May",
	time.June:      "Jun",
	time.July:      "Jul",
	time.August:    "Aug",
	time.September: "Sep",
	time.October:   "Oct",
	time.November:  "Nov",
	time.December:  "Dec",
}

// Month represents a calendar month. Internally stored as first day.
type Month time.Time

// NewMonth truncates t to the first day of its month (UTC).
func NewMonth(t time.Time) Month {
	return Month(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC))
}

// ParseMonth parses a "2026-05" string into a Month.
func ParseMonth(s string) (Month, error) {
	t, err := time.Parse("2006-01", s)
	if err != nil {
		return Month{}, ErrInvalidMonth
	}
	return NewMonth(t), nil
}

// String formats the month as "2026-05".
func (m Month) String() string {
	return time.Time(m).Format("2006-01")
}

// Label formats the month as "May 2026".
func (m Month) Label() string {
	return fmt.Sprintf("%s %d", monthAbbr[time.Time(m).Month()], time.Time(m).Year())
}

// Time returns the underlying time.Time (the first day of the month).
func (m Month) Time() time.Time {
	return time.Time(m)
}
