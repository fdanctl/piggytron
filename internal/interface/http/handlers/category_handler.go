package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	incomecategoryapp "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type CategoriesHandler struct {
	incomeCatService  *incomecategoryapp.Service
	expenseCatService *expensecategoryapp.Service
}

func NewCategoriesHandler(
	es *expensecategoryapp.Service,
	is *incomecategoryapp.Service,
) *CategoriesHandler {
	return &CategoriesHandler{
		incomeCatService:  is,
		expenseCatService: es,
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
		h.GetId(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	ec, err := h.expenseCatService.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	var ecView []views.ExpenseCategory
	for _, v := range ec {
		ecView = append(ecView, views.ExpenseCategory{
			Id:          v.ID(),
			Name:        v.Name(),
			ExpenseType: v.ExpenseType(),
		})
	}

	ic, err := h.incomeCatService.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	var icView []views.IncomeCategory
	for _, v := range ic {
		icView = append(icView, views.IncomeCategory{
			Id:   v.ID(),
			Name: v.Name(),
		})
	}

	content := partials.Categories(
		views.CategoriesView{
			IncomeCategories:  icView,
			ExpenseCategories: ecView,
		},
	)
	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
		return
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Categories").Render(ctx, w)
}

func (h *CategoriesHandler) GetId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	content := components.Breadcrumbs([]components.BreadcrumbsLinks{
		{
			Href: "/categories",
			Name: "Categories",
		},
		{
			Href: "/categories/" + id,
			Name: id,
		},
	}, []components.BreadcrumbsLinks{
		{
			Href: "/categories",
			Name: "hello",
		},
		{
			Href: "/categories",
			Name: "hello",
		},
	})

	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
		return
	}
	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Categories").Render(ctx, w)
}
