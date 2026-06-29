package handlers

import (
	"context"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

type FilterDialogHandler struct {
	categoryQueryService query.CategoryQueryService
	accountService       *appaccount.Service
	tQueryService        query.TransactionQueryService
	accQueryService      query.AccountQueryService
}

func NewFilterDialogHandler(
	cs query.CategoryQueryService,
	as *appaccount.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *FilterDialogHandler {
	return &FilterDialogHandler{
		categoryQueryService: cs,
		accountService:       as,
		tQueryService:        tq,
		accQueryService:      aq,
	}
}

func (h *FilterDialogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *FilterDialogHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	c, err := h.categoryQueryService.FindAllCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find all categories: %w", err))
		return
	}

	var categoryOptions []partials.FilterOption
	for _, v := range c {
		categoryOptions = append(
			categoryOptions,
			partials.FilterOption{Label: v.Name, Value: v.ID},
		)
	}

	account, err := h.accountService.FindAllByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	minA, maxA, minD, maxD, err := h.tQueryService.GetMinMaxAmountAndDate(
		r.Context(),
		sessionInfo.UserID,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	minA = int(math.Floor(float64(minA) / float64(100)))
	maxA = int(math.Ceil(float64(maxA) / float64(100)))
	if minA == maxA {
		maxA++
	}

	var accountOptions []partials.FilterOption
	for _, v := range account {
		accountOptions = append(
			accountOptions,
			partials.FilterOption{Label: v.Name(), Value: string(v.ID())},
		)
	}

	q := r.URL.Query()
	types := q["types"]
	accounts := q["accounts"]
	cats := q["categories"]
	minAmount := q.Get("minamount")
	maxAmount := q.Get("maxamount")
	minDate := q.Get("mindate")
	maxDate := q.Get("maxdate")

	filters := query.NewTransactionFilters(
		types,
		accounts,
		cats,
		minAmount,
		maxAmount,
		minDate,
		maxDate,
	)

	resCount, err := h.tQueryService.CountFilteredResults(
		r.Context(), sessionInfo.UserID, filters,
	)
	if err != nil {
		httperror.SendError(
			w,
			r,
			fmt.Errorf("failed to count filters transactions results: %w", err),
		)
		return
	}

	includedAcc, err := h.accQueryService.FindIDNamesIncludes(r.Context(), accounts)
	if err != nil {
		httperror.SendError(
			w,
			r,
			fmt.Errorf("failed to find accounts id-names for %v: %w", accounts, err),
		)
		return
	}
	includedCats, err := h.categoryQueryService.FindCategoriesIDIncludes(r.Context(), cats)
	if err != nil {
		httperror.SendError(
			w,
			r,
			fmt.Errorf("failed to find categories id-names for %v: %w", cats, err),
		)
		return
	}

	content := partials.TransactionsFilters(
		accountOptions,
		categoryOptions,
		includedAcc,
		includedCats,
		r.URL.Query(),
		minA, maxA,
		int(minD.Unix()), int(maxD.Unix()),
		resCount,
	)
	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("right-sheet", "Filters", nil).Render(ctx, w)
}

func (h *FilterDialogHandler) Post(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	r.ParseForm()
	q := r.Form
	types := q["types"]
	accounts := q["accounts"]
	cats := q["categories"]
	minAmount := q.Get("minamount")
	maxAmount := q.Get("maxamount")
	minDate := q.Get("mindate")
	maxDate := q.Get("maxdate")

	filters := query.NewTransactionFilters(
		types,
		accounts,
		cats,
		minAmount,
		maxAmount,
		minDate,
		maxDate,
	)

	filterCount, queries := queryStrFromFiltersWithCount(
		2,
		types,
		accounts,
		cats,
		minAmount,
		maxAmount,
		minDate,
		maxDate,
	)

	resCount, err := h.tQueryService.CountFilteredResults(
		r.Context(), sessionInfo.UserID, filters,
	)
	if err != nil {
		httperror.SendError(
			w,
			r,
			fmt.Errorf("failed to count filters transactions results: %w", err),
		)
		return
	}

	w.Header().Set("HX-Push-Url", "?"+strings.Join(queries[1:], "&"))
	w.Header().Set("HX-Trigger-After-Settle", "refetch-transactions")

	if filterCount > 0 {
		components.Button(
			"Reset",
			"w-full",
			components.BtnOutline,
			components.BtnMedium,
			templ.Attributes{
				"type":        "button",
				"data-action": "ui.filters.reset",
			},
		).Render(r.Context(), w)
	}
	components.Button(
		fmt.Sprintf("Show %d results", resCount),
		"w-full",
		components.BtnPrimary,
		components.BtnMedium,
		templ.Attributes{
			"type":        "button",
			"data-action": "ui.dialog.close-last",
		},
	).Render(r.Context(), w)

	templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p id=\"filter-result-count\" hx-swap-oob=\"innerHTML\">")
		if err != nil {
			return err
		}
		_, err = fmt.Fprintf(w, "%d results", resCount)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, "</p>")
		return err
	}).Render(r.Context(), w)
	components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
		"style":       "height: 24px;",
		"hx-swap-oob": "outerHTML",
		"id":          "filter-btn",
		"hx-get":      "/partials/transaction-filters?" + strings.Join(queries[1:], "&"),
		"hx-target":   "#dialog-root",
	}).Render(r.Context(), w)
}
