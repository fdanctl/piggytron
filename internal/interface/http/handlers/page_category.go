package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

// CategoriesHandler renders the categories overview and the per-category
// detail page.
type CategoriesHandler struct {
	categoryQuery      query.CategoryQueryService
	ledgerQueryService query.LedgerQueryService
}

func NewCategoriesHandler(
	cq query.CategoryQueryService,
	tq query.LedgerQueryService,
) *CategoriesHandler {
	return &CategoriesHandler{
		categoryQuery:      cq,
		ledgerQueryService: tq,
	}
}

func (h *CategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if id == "" {
			h.Get(w, r)
			return
		}
		h.GetWithID(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the categories overview, grouped into income and expense.
func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	categories, err := h.categoryQuery.FindAllCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	content := pages.Categories(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Categories"},
			},
			Options: nil,
		},
		categories,
	)

	renderWithMainLayout(w, r, "Categories", content)
}

// GetWithID renders one category (income or expense) with its paginated
// entries and sibling-category breadcrumb links.
func (h *CategoriesHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")

	category, err := h.categoryQuery.FindByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		httperror.SendError(w, r, err)
		return
	}

	categories, err := h.categoryQuery.FindAllCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var optionsLinks []views.BreadcrumbsLink

	for _, v := range categories {
		optionsLinks = append(optionsLinks, views.BreadcrumbsLink{
			Href: fmt.Sprintf("/categories/%s", v.ID),
			Name: v.Name,
		})
	}

	filters := query.NewLedgerFilters(nil, nil, []string{id}, "", "", "", "")

	transactions, err := h.ledgerQueryService.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*1-LIMIT,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("error finding filtered transaction: %w", err))
		return
	}
	var hasMore bool
	if len(transactions.Data) == LIMIT+1 {
		hasMore = true
		transactions.Data = transactions.Data[0 : len(transactions.Data)-1]
	}

	var transactionsView []views.Transaction
	for _, t := range transactions.Data {
		transactionsView = append(
			transactionsView,
			views.NewTransaction(t),
		)
	}

	content := pages.Category(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{
					Href: "/categories",
					Name: "Categories",
				},
				{
					Href: "/categories/" + category.ID,
					Name: category.Name,
				},
			},
			Options: optionsLinks,
		}, category, transactionsView, hasMore, transactions.Total,
	)

	renderWithMainLayout(w, r, category.Name, content)
}
