package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	incomecategoryapp "github.com/fdanctl/piggytron/internal/application/income_category"
	"github.com/fdanctl/piggytron/web/templates/components"
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
		// submit form

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *IncomeCategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	fmt.Println("income cat form")
	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := io.WriteString(w, "<p>in construction</p>")
		if err != nil {
			return err
		}
		err = components.Button(
			"logout",
			"w-fit",
			components.BtnDestructive,
			components.BtnMedium,
			templ.Attributes{
				"hx-get": "/partials/auth/logout",
			},
		).Render(ctx, w)
		return err
	})
	ctx := templ.WithChildren(r.Context(), content)
	components.DialogWrapper("", nil).Render(ctx, w)
}
