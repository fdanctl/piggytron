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
		if _, err := io.WriteString(w, "<div class=\"flex flex-wrap\">"); err != nil {
			return err
		}

		if err := components.CircleXSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.EyeSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.EyeOffSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.SearchSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.HouseSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ChartPieSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.PiggyBankSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.LandmarkSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.NotebookSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.TagSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.FileChartColumnSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.Settings2Svg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.SunSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.MoonSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ChevronLeftSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ChevronRightSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ChevronDownSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ChevronUpSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.CalendarDaysSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.MenuSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.XSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ArrowUpRight(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.PenSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.TrashSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.BellSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.PlusSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ArrowLeftSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ArrowRightSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.ArrowRightLeftSvg(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if err := components.Spinner(30, "", "").Render(ctx, w); err != nil {
			return err
		}

		if _, err := io.WriteString(w, "</div>"); err != nil {
			return err
		}

		if _, err := io.WriteString(w, "<div class=\"flex flex-col\">"); err != nil {
			return err
		}
		search := components.SearchBar("", nil)
		if err := search.Render(ctx, w); err != nil {
			return err
		}

		btn := components.Button(
			"logout",
			"w-fit",
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
			"current time",
			"w-fit",
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

		if err := components.Loader("", nil).Render(ctx, w); err != nil {
			return err
		}
		_, err := io.WriteString(w, "</div>")

		return err
	})

	if r.Header.Get("Hx-Request") == "true" {
		contents.Render(r.Context(), w)
		return
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, contents)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Dashboard").Render(ctx, w)
}

type BudgetHandler struct{}

func (h *BudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BudgetHandler) Get(w http.ResponseWriter, r *http.Request) {
	form := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
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
	if r.Header.Get("Hx-Request") == "true" {
		form.Render(r.Context(), w)
		return
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, form)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Budget").Render(ctx, w)
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
		ctx = templ.WithChildren(ctx, form)
		err := layouts.LogLayout().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("Login").Render(ctx, w)
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
		ctx = templ.WithChildren(ctx, form)
		err := layouts.LogLayout().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), layout)
	layouts.Base("Signup").Render(ctx, w)
}

type ExpensesHandler struct{}

func (h *ExpensesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ExpensesHandler) Get(w http.ResponseWriter, r *http.Request) {
	content := partials.Test()
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
	layouts.Base("Budget").Render(ctx, w)
}
