package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/a-h/templ"
	expensecategoryapp "github.com/fdanctl/piggytron/internal/application/expense_category"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type ExpenseCategoriesHandler struct {
	service *expensecategoryapp.Service
}

func NewExpenseCategoriesHandler(s *expensecategoryapp.Service) *ExpenseCategoriesHandler {
	return &ExpenseCategoriesHandler{
		service: s,
	}
}

func (h *ExpenseCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExpenseCategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := partials.ExpenseCategoryForm(*views.NewExpenseCategoryForm())
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", nil).Render(ctx, w)
}

func (h *ExpenseCategoriesHandler) Post(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	catType := r.FormValue("type")

	view := views.ExpenseCategoryForm{
		Initial: false,
		Name:    name,
		Type:    catType,
	}

	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.ExpenseCategoryForm(view).Render(r.Context(), w)
		return
	}

	category, err := h.service.CreateCategory(r.Context(), name, catType)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, expensecategoryapp.ErrDuplicate) {
			view.CustomError = err
		}
		partials.ExpenseCategoryForm(view).Render(r.Context(), w)
		return
	}

	ecView := views.ExpenseCategory{
		ID:          category.ID(),
		Name:        category.Name(),
		ExpenseType: category.ExpenseType(),
	}

	oob := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(
			w,
			"<ul hx-swap-oob=\"beforeend:#expense-cat ul\">",
		); err != nil {
			return err
		}

		if err := partials.CategoryItem(ecView, templ.Attributes{
			"style": "animation-delay: 0s;",
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, "</ul>")
		return err
	})

	w.Header().Set("HX-Trigger", "expenseCategoryAdded")
	partials.ExpenseCategoryForm(view).Render(r.Context(), w)
	oob.Render(r.Context(), w)
}
