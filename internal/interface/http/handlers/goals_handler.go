package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	transactionapp "github.com/fdanctl/piggytron/internal/application/transaction"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalsHandler struct {
	accountService      *accountapp.Service
	transactionService  *transactionapp.Service // will not be needed after TODO
	tQueryService       query.TransactionQueryService
	accountQueryService query.AccountQueryService
}

func NewGoalsHandler(
	ac *accountapp.Service,
	ts *transactionapp.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *GoalsHandler {
	return &GoalsHandler{
		accountService:      ac,
		transactionService:  ts,
		tQueryService:       tq,
		accountQueryService: aq,
	}
}

func (h *GoalsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if id == "" {
			h.Get(w, r)
			return
		}
		h.GetWithID(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (h *GoalsHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := sessionInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goals, err := h.accountQueryService.FindAllGoalsWithSum(
		r.Context(), sessionInfo.UserID,
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var gView []views.Goal
	for _, g := range goals {
		gView = append(gView, views.NewGoal(g))
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

func (h *GoalsHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := sessionInfoFromCtx(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := r.PathValue("id")
	goal, err := h.accountQueryService.FindOneWithSum(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	goals, err := h.accountQueryService.FindGoalsIDNames(r.Context(), sessionInfo.UserID)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var optionsLinks []components.BreadcrumbsLink
	for _, g := range goals {
		optionsLinks = append(optionsLinks, components.BreadcrumbsLink{
			Href: fmt.Sprintf("/goals/%s", g.ID),
			Name: g.Name,
		})
	}

	filters, err := query.NewTransactionFilters(nil, []string{id}, nil, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := h.tQueryService.FindFiltered(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*1-LIMIT,
	)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var hasMore bool
	if len(transactions) == LIMIT+1 {
		hasMore = true
		transactions = transactions[0 : len(transactions)-1]
	}

	var transactionsViews []views.Transaction
	for _, t := range transactions {
		transactionsViews = append(
			transactionsViews,
			views.NewAccountTransaction(t, &goal.AccountWithCategory),
		)
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		breadcrumbs := components.Breadcrumbs([]components.BreadcrumbsLink{
			{
				Href: "/goals",
				Name: "Goals",
			},
			{
				Href: "/goals/" + string(goal.ID),
				Name: goal.Name,
			},
		}, optionsLinks)
		if err := breadcrumbs.Render(ctx, w); err != nil {
			return err
		}

		return partials.Goal(views.NewGoal(goal), transactionsViews, hasMore).
			Render(ctx, w)
	})

	renderWithMainLayout(w, r, goal.Name, content)
}
