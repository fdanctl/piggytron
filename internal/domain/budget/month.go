package budget

import "time"

// Month represents a calendar month. Internally stored as first day.
type Month time.Time

func NewMonth(t time.Time) Month {
	return Month(time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC))
}

func ParseMonth(s string) (Month, error) {
	t, err := time.Parse("2006-01", s) // parses "2026-05"
	if err != nil {
		return Month{}, ErrInvalidMonth
	}
	return NewMonth(t), nil
}

func (m Month) String() string {
	return time.Time(m).Format("2006-01") // outputs "2026-05"
}

func (m Month) Time() time.Time {
	return time.Time(m)
}
