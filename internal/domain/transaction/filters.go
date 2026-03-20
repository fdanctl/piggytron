package transaction

import (
	"strconv"
)

type Filters struct {
	Ttypes *[]Ttype

	AccountIds *[]ID

	CategoryIds *[]ID

	MinAmount *uint
	MaxAmount *uint

	// MinDate *time.Time
	// MaxDate *time.Time
}

func NewFilters(
	ttype, accountIds, categoryIds []string,
	minAmount, maxAmount string,
	// minDate, maxDate string,
) (*Filters, error) {
	var pt *[]Ttype
	var pa *[]ID
	var pc *[]ID
	var pmin *uint
	var pmax *uint

	if len(ttype) > 0 {
		var tts []Ttype
		for _, t := range ttype {
			ttype, err := NewType(t)
			if err != nil {
				return nil, err
			}
			tts = append(tts, ttype)
		}
		pt = &tts
	}

	if len(accountIds) > 0 {
		var accIds []ID
		for _, a := range accountIds {
			acc, err := NewId(a)
			if err != nil {
				return nil, err
			}
			accIds = append(accIds, acc)
		}
		pa = &accIds
	}

	if len(categoryIds) > 0 {
		var catIds []ID
		for _, c := range categoryIds {
			cat, err := NewId(c)
			if err != nil {
				return nil, err
			}
			catIds = append(catIds, cat)
		}
		pc = &catIds
	}

	if minAmount != "" {
		minAmnt, err := strconv.Atoi(minAmount)
		if err != nil {
			return nil, ErrNotNumber
		}
		if minAmnt < 0 {
			return nil, ErrNegativeNumber
		}
		mina := uint(minAmnt)
		pmin = &mina
	}

	if maxAmount != "" {
		maxAmnt, err := strconv.Atoi(maxAmount)
		if err != nil {
			return nil, ErrNotNumber
		}
		if maxAmnt < 0 {
			return nil, ErrNegativeNumber
		}
		maxa := uint(maxAmnt)
		pmax = &maxa
	}

	if (pmin != nil && pmax != nil) && *pmax < *pmin {
		return nil, ErrMinMax
	}

	return &Filters{
		Ttypes:      pt,
		AccountIds:  pa,
		CategoryIds: pc,
		MinAmount:   pmin,
		MaxAmount:   pmax,
	}, nil
}
