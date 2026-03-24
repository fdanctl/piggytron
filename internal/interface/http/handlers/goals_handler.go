package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalsHandler struct {
	accountService     *accountapp.Service
	transactionService *transactionapp.Service
}

func NewGoalsHandler(ac *accountapp.Service, ts *transactionapp.Service) *GoalsHandler {
	return &GoalsHandler{
		accountService:     ac,
		transactionService: ts,
	}
}

func (h *GoalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *GoalsHandler) Get(w http.ResponseWriter, r *http.Request) {
	goals, err := h.accountService.ReadAllGoalsByUser(r.Context())
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var gView []views.Goal
	for _, g := range goals {
		gView = append(gView, views.NewGoal(g, 0))
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Goals"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		err = partials.Goals(gView).Render(ctx, w)
		return err
	})

	if r.Header.Get("Hx-Request") == "true" {
		content.Render(r.Context(), w)
		io.WriteString(w, "<title>Goals</title>")
		return
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	layouts.Base("Goals").Render(ctx, w)
}
