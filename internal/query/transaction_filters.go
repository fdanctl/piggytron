package query

import (
	"errors"
	"strconv"
)

var (
	ErrNotNumber      = errors.New("not a number")
	ErrNegativeNumber = errors.New("number can't be negative")
	ErrMinMax         = errors.New("max in lower than min")
)

type TransactionFilters struct {
	Types *[]string

	AccountIds *[]string

	CategoryIds *[]string

	MinAmount *int
	MaxAmount *int

	// MinDate *time.Time
	// MaxDate *time.Time
}

func NewTransactionFilters(
	ttype, accountIds, categoryIds []string,
	minAmount, maxAmount string,
	// minDate, maxDate string,
) (*TransactionFilters, error) {
	// var pt *[]string
	// var pa *[]string
	// var pc *[]string
	var pmin *int
	var pmax *int

	// if len(ttype) > 0 {
	// 	var tts []string
	// 	for _, t := range ttype {
	// 		tts = append(tts, ttype)
	// 	}
	// 	pt = &tts
	// }
	//
	// if len(accountIds) > 0 {
	// 	var accIds []transaction.ID
	// 	for _, a := range accountIds {
	// 		acc, err := transaction.NewId(a)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		accIds = append(accIds, acc)
	// 	}
	// 	pa = &accIds
	// }
	//
	// if len(categoryIds) > 0 {
	// 	var catIds []transaction.ID
	// 	for _, c := range categoryIds {
	// 		cat, err := transaction.NewId(c)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		catIds = append(catIds, cat)
	// 	}
	// 	pc = &catIds
	// }

	if minAmount != "" {
		minAmnt, err := strconv.Atoi(minAmount)
		if err != nil {
			return nil, ErrNotNumber
		}
		if minAmnt < 0 {
			return nil, ErrNegativeNumber
		}
		mina := int(minAmnt)
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
		maxa := int(maxAmnt)
		pmax = &maxa
	}

	if (pmin != nil && pmax != nil) && *pmax < *pmin {
		return nil, ErrMinMax
	}

	return &TransactionFilters{
		Types:       &ttype,
		AccountIds:  &accountIds,
		CategoryIds: &categoryIds,
		MinAmount:   pmin,
		MaxAmount:   pmax,
	}, nil
}
