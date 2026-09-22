package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/templates/partials"
)

// FilterDialogHandler renders the ledger filter dialog (GET) and applies
// the submitted filters (POST).
type FilterDialogHandler struct {
	categoryQueryService query.CategoryQueryService
	accountService       *appaccount.Service
	tQueryService        query.LedgerQueryService
	accQueryService      query.AccountQueryService
}

func NewFilterDialogHandler(
	cs query.CategoryQueryService,
	as *appaccount.Service,
	tq query.LedgerQueryService,
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

// Get renders the filter form: category/account options, amount and date
// bounds from the user's entries, and the currently applied filters.
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

	filters := query.NewLedgerFilters(
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
			fmt.Errorf("failed to count filtered ledger entries results: %w", err),
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

	content := partials.LedgerFilters(
		accountOptions,
		categoryOptions,
		includedAcc,
		includedCats,
		r.URL.Query(),
		minA, maxA,
		int(minD.Unix()), int(maxD.Unix()),
	)
	bottom := partials.LedgerFiltersBtns(resCount)

	components.DialogWrapperForm(
		"dialog--right-sheet",
		components.DialogHeader("", "Filters", nil),
		content,
		bottom,
		templ.Attributes{"id": "transactions-filters"},
	).Render(r.Context(), w)
}

// Post applies the submitted filters: it pushes the filter query string to
// the URL, triggers a transaction refetch and updates the result count and
// the filter button badge.
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

	filters := query.NewLedgerFilters(
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
			fmt.Errorf("failed to count filtered ledger entries results: %w", err),
		)
		return
	}

	w.Header().Set("HX-Push-Url", "?"+strings.Join(queries[1:], "&"))
	w.Header().Set("HX-Trigger", "refetch-transactions")

	templ.Join(
		partials.LedgerFiltersBtns(resCount),
		layouts.HxPartial(
			"#filter-result-count",
			"outerHTML",
			pages.FilterResultCount(resCount),
		),
		layouts.HxPartial(
			"#filter-btn",
			"outerHTML",
			components.FilterBtn(uint8(filterCount), 0, "", "", templ.Attributes{
				"style":     "height: 24px;",
				"hx-get":    "/partials/ledger-filters?" + strings.Join(queries[1:], "&"),
				"hx-target": "#dialog-root",
			}),
		),
	).Render(r.Context(), w)
}
