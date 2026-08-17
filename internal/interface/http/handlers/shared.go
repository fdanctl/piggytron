// Package handlers implements the HTTP layer: page handlers (full HTML) and
// partial handlers (HTMX fragments, including the chart partials), plus the
// helpers they share. Handlers depend only on application services and query
// interfaces — never on Postgres directly.
package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/domain/account"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

var (
	// ErrInvalidAmount signals an amount string that cannot be parsed into
	// cents.
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrEmpty signals a required string argument that is empty.
	ErrEmpty = errors.New("string is empty")
)

// LIMIT is the default page size for paginated ledger listings; handlers
// request LIMIT+1 rows to detect whether more pages exist.
const LIMIT = 30

func renderWithMainLayout(
	w http.ResponseWriter,
	r *http.Request,
	title string,
	content templ.Component,
) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("Hx-Request") == "true" {
		err := content.Render(r.Context(), w)
		fmt.Fprintf(w, "<title>%s</title>", title)
		return err
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	return layouts.Base(title).Render(ctx, w)
}

func convertAmountStrToInt(str string) (int, error) {
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
			return 0, ErrInvalidAmount
		}
		return parsed * 100, nil
	}

	if length-1-i > 2 {
		return 0, ErrInvalidAmount
	}

	for length-i < 3 {
		str += "0"
		length++
	}

	tAmount, err := strconv.Atoi(strings.Replace(str, ".", "", 1))
	if err != nil {
		return 0, ErrInvalidAmount
	}

	return tAmount, nil
}

// parseMonth receive a string of type 042026,
// and return the correspondent year and month
// TODO remove
func parseMonth(str string) (int, time.Month, error) {
	if len(str) != 6 {
		return 0, time.January, errors.New("wrong month")
	}

	m, err := strconv.Atoi(str[:2])
	if err != nil {
		return 0, time.January, errors.New("wrong month")
	}

	y, err := strconv.Atoi(str[2:])
	if err != nil {
		return 0, time.January, errors.New("wrong month")
	}
	return y, time.Month(m), nil
}

func queryStrFromFiltersWithCount(
	page int,
	types, accounts, cats []string,
	minAmount, maxAmount string,
	minDate, maxDate string,
) (int, []string) {
	filterCount := len(types) + len(accounts) + len(cats)
	queries := []string{fmt.Sprintf("page=%d", page)}
	if len(types) > 0 {
		queries = append(queries, "types="+strings.Join(types, "&types="))
	}
	if len(accounts) > 0 {
		queries = append(queries, "accounts="+strings.Join(accounts, "&accounts="))
	}
	if len(cats) > 0 {
		queries = append(queries, "categories="+strings.Join(cats, "&categories="))
	}
	if minAmount != "" || maxAmount != "" {
		if minAmount != "" {
			queries = append(queries, "minamount="+minAmount)
		}
		if maxAmount != "" {
			queries = append(queries, "maxamount="+maxAmount)
		}
		filterCount++
	}

	if minDate != "" || maxDate != "" {
		if minDate != "" {
			queries = append(queries, "mindate="+minDate)
		}
		if maxDate != "" {
			queries = append(queries, "maxdate="+maxDate)
		}
		filterCount++
	}
	return filterCount, queries
}

func getCategorySelectOptions(
	qs query.CategoryQueryService,
	ctx context.Context,
	userID string,
) (iCatOpts, eCatOpts []components.SelectOption, err error) {
	cats, err := qs.FindAllCategories(ctx, userID)
	if err != nil {
		return
	}
	for _, v := range cats {
		if v.Status == "archived" {
			continue
		}
		if v.Type == "income" {
			iCatOpts = append(
				iCatOpts,
				components.SelectOption{Label: v.Name, Value: v.ID},
			)
		} else {
			eCatOpts = append(
				eCatOpts,
				components.SelectOption{Label: v.Name, Value: v.ID},
			)
		}
	}

	return
}

func getAccSelectOptions(
	as *appaccount.Service,
	ctx context.Context,
	userID string,
) (noSavingsBanksOpts, goalsSavingsOpts []components.SelectOption, err error) {
	acc, err := as.FindAllByUser(ctx, userID)
	if err != nil {
		return
	}

	for _, v := range acc {
		if v.Status() == account.ClosedStatus {
			continue
		}
		if v.Type() == account.CheckingType {
			noSavingsBanksOpts = append(
				noSavingsBanksOpts,
				components.SelectOption{Label: v.Name(), Value: string(v.ID())},
			)
		} else {
			goalsSavingsOpts = append(
				goalsSavingsOpts,
				components.SelectOption{Label: v.Name(), Value: string(v.ID())},
			)
		}
	}
	return
}

func getStartPeriodDate(period string, clamp time.Time) time.Time {
	now := time.Now()
	switch period {
	case "all":
		return clamp
	case "1y":
		yearago := time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, now.Location())
		if clamp.Before(yearago) {
			return yearago
		}
	case "ytd":
		ytd := time.Date(now.Year(), time.January, 1, 0, 0, 0, 0, now.Location())
		if clamp.Before(ytd) {
			return ytd
		}
	case "6m":
		monthsago := time.Date(now.Year(), now.Month()-6, 1, 0, 0, 0, 0, now.Location())
		if clamp.Before(monthsago) {
			return monthsago
		}
	case "mtd":
		mtd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		if clamp.Before(mtd) {
			return mtd
		}
	}
	return clamp
}
