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

// LoginHandler serves the login page (and the HTMX login form fragment).
type LoginHandler struct {
	isDev bool
}

// NewLoginHandler builds the handler; in dev mode the form is pre-filled
// with demo credentials.
func NewLoginHandler(isDev bool) *LoginHandler {
	return &LoginHandler{
		isDev: isDev,
	}
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the login form, standalone or wrapped in the base layout.
func (h *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	v := views.NewLoginView(redirect)
	if h.isDev {
		v.Name = "gopher"
		v.Password = "123"
	}
	form := partials.LoginForm(*v)

	if r.Header.Get("Hx-Request") == "true" {
		form.Render(r.Context(), w)
		io.WriteString(w, "<title>Login</title>")
		return
	}

	layout := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, form)
		err := layouts.LogLayout().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("Login").Render(ctx, w)
}
