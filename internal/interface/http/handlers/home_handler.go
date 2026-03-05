package handlers

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *HomeHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	layouts.Base(
		"title",
		components.EyeSvg(80, ""),
		components.CircleXSvg(80, ""),
		components.Button(
			"logout",
			"",
			components.BtnDestructive,
			components.BtnMedium,
			templ.Attributes{
				"hx-get": "/partials/auth/logout",
			},
		),
	).Render(r.Context(), w)
}

type LoginHandler struct{}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *LoginHandler) Get(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	form := partials.LoginForm(*views.NewLoginView(redirect))
	if r.Header.Get("Hx-Request") == "true" {
		form.Render(r.Context(), w)
		return
	}
	layout := layouts.LogLayout(form)
	layouts.LogLayout()
	layouts.Base("login", layout).Render(r.Context(), w)
}

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
		return
	}
	layout := layouts.LogLayout(form)
	layouts.LogLayout()
	layouts.Base("login", layout).Render(r.Context(), w)
}
