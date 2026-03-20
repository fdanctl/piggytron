package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type AllTransactionsHandler struct {
	service *transactionapp.Service
}

func NewAllTransactionsHandler(s *transactionapp.Service) *AllTransactionsHandler {
	return &AllTransactionsHandler{
		service: s,
	}
}

func (h *AllTransactionsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *AllTransactionsHandler) Get(w http.ResponseWriter, r *http.Request) {
	// TODO infinite scroll (maybe 50 each time)
	transactions, err := h.service.ReadAllByUser(r.Context(), 1)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var transactionsView []views.Transaction
	for _, t := range transactions {
		transactionsView = append(
			transactionsView,
			views.NewTransaction(t.ID(), t.Description(), t.Ttype(), t.Amount(), t.Date()),
		)
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Transactions"},
			{Href: "", Name: "All"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		return partials.TransactionsList(transactionsView, "page=2").Render(ctx, w)
	})
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
	layouts.Base("Transactions").Render(ctx, w)
}
