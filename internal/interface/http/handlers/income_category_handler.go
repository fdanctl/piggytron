package handlers

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/a-h/templ"
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
	sessionInfo, err := sessionInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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

	category, err := h.service.CreateCategory(r.Context(), sessionInfo.UserID, name)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, incomecategoryapp.ErrDuplicate) {
			view.CustomError = err
		}
		partials.IncomeCategoryForm(view).Render(r.Context(), w)
		return
	}

	icView := views.IncomeCategory{
		ID:   category.ID(),
		Name: category.Name(),
	}

	oob := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if _, err := io.WriteString(
			w,
			"<ul hx-swap-oob=\"beforeend:#income-cat ul\">",
		); err != nil {
			return err
		}

		if err := partials.CategoryItem(icView, templ.Attributes{
			"style": "animation-delay: 0s;",
		}).Render(ctx, w); err != nil {
			return err
		}

		_, err := io.WriteString(w, "</ul>")
		return err
	})

	w.Header().Set("HX-Trigger", "incomeCategoryAdded")
	partials.IncomeCategoryForm(view).Render(r.Context(), w)
	oob.Render(r.Context(), w)
}
