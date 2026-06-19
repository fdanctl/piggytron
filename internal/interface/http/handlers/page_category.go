package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appexpensecategory"
	"github.com/fdanctl/piggytron/internal/application/appincomecategory"
	"github.com/fdanctl/piggytron/internal/domain/incomecategory"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type CategoriesHandler struct {
	incomeCatService  *appincomecategory.Service
	expenseCatService *appexpensecategory.Service
	tQueryService     query.TransactionQueryService
}

func NewCategoriesHandler(
	es *appexpensecategory.Service,
	is *appincomecategory.Service,
	tq query.TransactionQueryService,
) *CategoriesHandler {
	return &CategoriesHandler{
		incomeCatService:  is,
		expenseCatService: es,
		tQueryService:     tq,
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

func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	// TODO make it one query intead of two
	ec, err := h.expenseCatService.FindAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	var ecView []views.ExpenseCategory
	for _, v := range ec {
		ecView = append(ecView, views.NewExpenseCategory(v))
	}

	ic, err := h.incomeCatService.FindAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var icView []views.IncomeCategory
	for _, v := range ic {
		icView = append(icView, views.NewIncomeCategory(v))
	}

	content := pages.Categories(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Categories"},
			},
			Options: nil,
		},
		views.CategoriesView{
			IncomeCategories:  icView,
			ExpenseCategories: ecView,
		},
	)

	renderWithMainLayout(w, r, "Categories", content)
}

func (h *CategoriesHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	id := r.PathValue("id")
	ecat, err := h.expenseCatService.FindCategory(r.Context(), id, sessionInfo.UserID)
	var icat *incomecategory.IncomeCategory
	if err != nil {
		icat, err = h.incomeCatService.FindCategory(r.Context(), id, sessionInfo.UserID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.NotFound(w, r)
				return
			}
			httperror.SendError(w, r, err)
			return
		}
	}

	// TODO use category query to be one query instead of two
	icats, err := h.incomeCatService.FindAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	ecats, err := h.expenseCatService.FindAllUserCategories(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var category views.Category
	if ecat != nil {
		category = views.NewExpenseCategory(ecat)
	} else {
		category = views.NewIncomeCategory(icat)
	}

	var optionsLinks []views.BreadcrumbsLink

	for _, v := range icats {
		optionsLinks = append(optionsLinks, views.BreadcrumbsLink{
			Href: fmt.Sprintf("/categories/%s", v.ID()),
			Name: v.Name(),
		})
	}
	for _, v := range ecats {
		optionsLinks = append(optionsLinks, views.BreadcrumbsLink{
			Href: fmt.Sprintf("/categories/%s", v.ID()),
			Name: v.Name(),
		})
	}

	filters := query.NewTransactionFilters(nil, nil, []string{id}, "", "", "", "")

	transactions, err := h.tQueryService.FindFilteredWithCount(
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
					Href: "/categories/" + category.GetID(),
					Name: category.GetName(),
				},
			},
			Options: optionsLinks,
		}, category, transactionsView, hasMore, transactions.Total,
	)

	renderWithMainLayout(w, r, category.GetName(), content)
}
