package query

import (
	"strconv"
	"time"

	"github.com/google/uuid"
)

// LedgerFilters is an optional filter set for ledger queries. Nil fields
// mean "no filter"; non-nil slices/pointers are treated as an OR set (any
// match is included).
type LedgerFilters struct {
	Types *[]string

	AccountIDs *[]string

	CategoryIDs *[]string

	MinAmount *int
	MaxAmount *int

	MinDate *time.Time
	MaxDate *time.Time
}

// NewLedgerFilters builds LedgerFilters from raw HTTP query strings,
// discarding invalid values. Amounts are given in the UI unit (e.g. 12.50 as
// "12.50") and converted to cents; dates are Unix seconds. An invalid or
// contradictory range yields nil for the affected fields.
func NewLedgerFilters(
	ttype, accountIDs, categoryIDs []string,
	minAmount, maxAmount string,
	minDate, maxDate string,
) *LedgerFilters {
	var ptypes *[]string
	tt := make([]string, 0, len(ttype))
	var paccs *[]string
	accs := make([]string, 0, len(accountIDs))
	var pcats *[]string
	cats := make([]string, 0, len(categoryIDs))
	var pminA *int
	var pmaxA *int
	var pminD *time.Time
	var pmaxD *time.Time

	for _, t := range ttype {
		if t == "income" || t == "expense" || t == "transfer" || t == "interest" ||
			t == "goal-fulfillment" ||
			t == "initial-balance" {
			tt = append(tt, t)
		}
	}
	if len(tt) > 0 {
		ptypes = &tt
	}

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
		minAmnt, err := strconv.Atoi(minAmount + "00")
		if err != nil || minAmnt < 0 {
			pminA = nil
		} else {
			mina := int(minAmnt)
			pminA = &mina
		}
	}

	if maxAmount != "" {
		maxAmnt, err := strconv.Atoi(maxAmount + "00")
		if err != nil || maxAmnt < 0 {
			pmaxA = nil
		} else {
			maxa := int(maxAmnt)
			pmaxA = &maxa
		}
	}

	if (pminA != nil && pmaxA != nil) && *pmaxA < *pminA {
		pmaxA = nil
		pminA = nil
	}

	if minDate != "" {
		minD, err := strconv.ParseInt(minDate, 10, 64)
		if err != nil || minD < 0 {
			pminD = nil
		} else {
			mind := time.Unix(minD, 0)
			// make it the start of the day for the db include this date
			mind = time.Date(mind.Year(), mind.Month(), mind.Day(), 0, 0, 0, 0, mind.Location())
			pminD = &mind
		}
	}

	if maxDate != "" {
		maxD, err := strconv.ParseInt(maxDate, 10, 64)
		if err != nil || maxD < 0 {
			pmaxD = nil
		} else {
			maxd := time.Unix(maxD, 0)
			// don't make it the start of the day for the db include this date
			pmaxD = &maxd
		}
	}

	if (pminD != nil && pmaxD != nil) && pmaxD.Compare(*pminD) < 0 {
		pmaxD = nil
		pminD = nil
	}

	return &LedgerFilters{
		Types:       ptypes,
		AccountIDs:  paccs,
		CategoryIDs: pcats,
		MinAmount:   pminA,
		MaxAmount:   pmaxA,
		MinDate:     pminD,
		MaxDate:     pmaxD,
	}
}
