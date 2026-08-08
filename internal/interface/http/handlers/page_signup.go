package handlers

import (
	"context"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type SignupHandler struct{}

func (h *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *SignupHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := partials.SignupForm(*views.NewSignupView())
	if r.Header.Get("Hx-Request") == "true" {
		form.Render(r.Context(), w)
		io.WriteString(w, "<title>Signup</title>")
		return
	}

	layout := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, form)
		err := layouts.LogLayout().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("Signup").Render(ctx, w)
}
