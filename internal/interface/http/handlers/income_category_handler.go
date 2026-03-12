package handlers

import (
	"errors"
	"net/http"

	"github.com/a-h/templ"
	incomecategory "github.com/fdanctl/piggytron/internal/application/income_category"
	incomecategoryapp "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type IncomeCategoriesHandler struct {
	service *incomecategoryapp.Service
}

func NewIncomeCategoriesHandler(s *incomecategoryapp.Service) *IncomeCategoriesHandler {
	return &IncomeCategoriesHandler{
		service: s,
	}
}

func (h *IncomeCategoriesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *IncomeCategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := partials.IncomeCategoryForm(*views.NewIncomeCategoryForm())
	ctx := templ.WithChildren(r.Context(), form)
	components.DialogWrapper("", nil).Render(ctx, w)
}

func (h *IncomeCategoriesHandler) Post(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	view := views.IncomeCategoryForm{
		Initial: false,
		Name:    name,
	}

	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.IncomeCategoryForm(view).Render(r.Context(), w)
		return
	}

	err := h.service.CreateCategory(r.Context(), name)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, incomecategory.ErrDuplicate) {
			view.CustomError = err
		}
		partials.IncomeCategoryForm(view).Render(r.Context(), w)
		return
	}

	ic, err := h.service.ReadAllUserCategories(r.Context())
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	var icView []views.IncomeCategories
	for _, v := range ic {
		icView = append(icView, views.IncomeCategories{
			Id:   v.ID(),
			Name: v.Name(),
		})
	}

	w.Header().Set("HX-Trigger", "incomeCategoryAdded")

	partials.IncomeCategoryForm(view).Render(r.Context(), w)
	partials.IncomeCategories(icView, templ.Attributes{
		"hx-swap-oob": "outerHTML:#income-cat ul",
	}).Render(r.Context(), w)
}
