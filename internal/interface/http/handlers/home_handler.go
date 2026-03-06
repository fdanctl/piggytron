package handlers

import (
	"context"
	"io"
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

	contents := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		if err := components.EyeSvg(80, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.CircleXSvg(80, "", "").Render(ctx, w); err != nil {
			return err
		}

		btn := components.Button(
			"logout",
			"",
			components.BtnDestructive,
			components.BtnMedium,
			templ.Attributes{
				"hx-get": "/partials/auth/logout",
			},
		)
		if err := btn.Render(ctx, w); err != nil {
			return err
		}

		btn2 := components.ButtonWithIcon(
			"corrent time",
			"",
			components.BtnOutline,
			components.BtnMedium,
			templ.Attributes{
				"hx-get": "/partials/slow",
			},
			components.IconLeft,
			components.Spinner(0, "", "indicator"),
		)
		if err := btn2.Render(ctx, w); err != nil {
			return err
		}
		return nil
	})

	ctx := templ.WithChildren(r.Context(), contents)
	layouts.Base("title").Render(ctx, w)
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

	layout := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(r.Context(), form)
		err := layouts.LogLayout().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("title").Render(ctx, w)
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

	layout := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(r.Context(), form)
		err := layouts.LogLayout().Render(ctx, w)
		ctx = templ.ClearChildren(ctx)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("title").Render(ctx, w)
}
