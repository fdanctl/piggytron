package handlers

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type GoalsHandler struct {
	accountService      *appaccount.Service
	tQueryService       query.TransactionQueryService
	accountQueryService query.AccountQueryService
}

func NewGoalsHandler(
	ac *appaccount.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *GoalsHandler {
	return &GoalsHandler{
		accountService:      ac,
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
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	goals, err := h.accountQueryService.FindAllGoalsWithSum(
		r.Context(), sessionInfo.UserID,
	)
	if err != nil {
		logger.Error("error query goals", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var gView []views.Goal
	for _, g := range goals {
		gView = append(gView, views.NewGoal(g))
	}

	content := pages.Goals(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Goals"},
			},
			Options: nil,
		},
		gView,
	)

	renderWithMainLayout(w, r, "Goals", content)
}

func (h *GoalsHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id := r.PathValue("id")
	goal, err := h.accountQueryService.FindWithSum(r.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.NotFound(w, r)
			return
		}
		logger.Error("error finding one account", "error", err, "aid", id)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	goals, err := h.accountQueryService.FindGoalsIDNames(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error query goals", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	var optionsLinks []views.BreadcrumbsLink
	for _, g := range goals {
		optionsLinks = append(optionsLinks, views.BreadcrumbsLink{
			Href: fmt.Sprintf("/goals/%s", g.ID),
			Name: g.Name,
		})
	}

	filters := query.NewTransactionFilters(nil, []string{id}, nil, "", "", "", "")

	transactions, err := h.tQueryService.FindFilteredWithCount(
		r.Context(),
		sessionInfo.UserID,
		filters,
		LIMIT+1,
		LIMIT*1-LIMIT,
	)
	if err != nil {
		logger.Error("error find filtered", "error", err, "filters", filters)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	var hasMore bool
	if len(transactions.Data) == LIMIT+1 {
		hasMore = true
		transactions.Data = transactions.Data[0 : len(transactions.Data)-1]
	}

	var transactionsViews []views.Transaction
	for _, t := range transactions.Data {
		transactionsViews = append(
			transactionsViews,
			views.NewAccountTransaction(t, goal.Name),
		)
	}

	content := pages.Goal(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{
					Href: "/goals",
					Name: "Goals",
				},
				{
					Href: "/goals/" + string(goal.ID),
					Name: goal.Name,
				},
			},
			Options: optionsLinks,
		},
		views.NewGoal(*goal), transactionsViews, hasMore, transactions.Total,
	)

	renderWithMainLayout(w, r, goal.Name, content)
}
