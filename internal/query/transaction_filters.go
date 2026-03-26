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

	AccountIDs *[]string

	CategoryIDs *[]string

	MinAmount *int
	MaxAmount *int

	// MinDate *time.Time
	// MaxDate *time.Time
}

func NewTransactionFilters(
	ttype, accountIDs, categoryIDs []string,
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
	// if len(accountIDs) > 0 {
	// 	var accIDs []transaction.ID
	// 	for _, a := range accountIDs {
	// 		acc, err := transaction.NewID(a)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		accIDs = append(accIDs, acc)
	// 	}
	// 	pa = &accIDs
	// }
	//
	// if len(categoryIDs) > 0 {
	// 	var catIDs []transaction.ID
	// 	for _, c := range categoryIDs {
	// 		cat, err := transaction.NewID(c)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		catIDs = append(catIDs, cat)
	// 	}
	// 	pc = &catIDs
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
		AccountIDs:  &accountIDs,
		CategoryIDs: &categoryIDs,
		MinAmount:   pmin,
		MaxAmount:   pmax,
	}, nil
}
