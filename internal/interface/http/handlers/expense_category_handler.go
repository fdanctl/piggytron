package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/a-h/templ"
	expensecategory "github.com/fdanctl/piggytron/internal/application/expense_category"
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
	catType, err := strconv.Atoi(r.FormValue("type"))
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.ExpenseCategoryForm(*views.NewExpenseCategoryForm()).Render(r.Context(), w)
		return
	}

	view := views.ExpenseCategoryForm{
		Initial: false,
		Name:    name,
		Type:    int8(catType),
	}

	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.ExpenseCategoryForm(view).Render(r.Context(), w)
		return
	}

	err = h.service.CreateCategory(r.Context(), name, catType)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, expensecategory.ErrDuplicate) {
			view.CustomError = err
		}
		partials.ExpenseCategoryForm(view).Render(r.Context(), w)
		return
	}

	ic, err := h.service.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	var icView []views.ExpenseCategories
	for _, v := range ic {
		icView = append(icView, views.ExpenseCategories{
			Id:          v.ID(),
			Name:        v.Name(),
			ExpenseType: v.ExpenseType(),
		})
	}

	w.Header().Set("HX-Trigger", "expenseCategoryAdded")

	partials.ExpenseCategoryForm(view).Render(r.Context(), w)
	partials.ExpenseCategories(icView, templ.Attributes{
		"hx-swap-oob": "outerHTML:#expense-cat ul",
	}).Render(r.Context(), w)
}
