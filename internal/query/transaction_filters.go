package query

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

type TransactionFilters struct {
	Types *[]string

	AccountIDs *[]string

	CategoryIDs *[]string

	MinAmount *int
	MaxAmount *int

	MinDate *time.Time
	MaxDate *time.Time
}

func NewTransactionFilters(
	ttype, accountIDs, categoryIDs []string,
	minAmount, maxAmount string,
	minDate, maxDate string,
) *TransactionFilters {
	var paccs *[]string
	accs := make([]string, 0, len(accountIDs))
	var pcats *[]string
	cats := make([]string, 0, len(categoryIDs))
	var pmin *int
	var pmax *int

	// TODO very very type

	for _, aid := range accountIDs {
		if _, err := uuid.Parse(aid); err == nil {
			accs = append(accs, aid)
		}
	}
	if len(accs) > 0 {
		paccs = &accs
	}

	for _, cid := range categoryIDs {
		if _, err := uuid.Parse(cid); err == nil {
			cats = append(cats, cid)
		}
	}
	if len(cats) > 0 {
		pcats = &cats
	}

	if minAmount != "" {
		minAmnt, err := strconv.Atoi(minAmount)
		if err != nil || minAmnt < 0 {
			pmin = nil
		} else {
			mina := int(minAmnt)
			pmin = &mina
		}
	}

	if maxAmount != "" {
		maxAmnt, err := strconv.Atoi(maxAmount)
		if err != nil || maxAmnt < 0 {
			pmax = nil
		} else {
			maxa := int(maxAmnt)
			pmax = &maxa
		}
	}

	if (pmin != nil && pmax != nil) && *pmax < *pmin {
		pmax = nil
		pmin = nil
	}

	return &TransactionFilters{
		Types:       &ttype,
		AccountIDs:  paccs,
		CategoryIDs: pcats,
		MinAmount:   pmin,
		MaxAmount:   pmax,
	}
}
