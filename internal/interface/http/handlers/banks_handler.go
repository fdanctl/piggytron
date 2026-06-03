package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	accountapp "github.com/fdanctl/piggytron/internal/application/account"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type BanksHandler struct {
	service          *accountapp.Service
	transactionQuery query.TransactionQueryService
	accountQuery     query.AccountQueryService
}

func NewBanksHandler(
	s *accountapp.Service,
	tq query.TransactionQueryService,
	aq query.AccountQueryService,
) *BanksHandler {
	return &BanksHandler{
		service:          s,
		transactionQuery: tq,
		accountQuery:     aq,
	}
}

func (h *BanksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		id := r.PathValue("id")
		if id == "" {
			h.Get(w, r)
			return
		}
		h.GetWithID(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *BanksHandler) Get(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	transactions, err := h.transactionQuery.GetRecentTransactions(
		r.Context(),
		sessionInfo.UserID,
		5,
	)
	if err != nil {
		logger.Error("error finding transactions", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error finding accounts", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	pageView := views.NewBankPage(accounts, transactions)

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Banks"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		err = partials.Banks(pageView).Render(ctx, w)
		return err
	})

	renderWithMainLayout(w, r, "Banks", content)
}

func (h *BanksHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		logger.Error("unexpected error", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	aid := r.PathValue("id")
	bank, err := h.accountQuery.FindWithSum(r.Context(), aid)
	if err != nil {
		logger.Error("error read one bank", "error", err, "account_id", aid)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	banks, err := h.accountQuery.FindBanksIDNames(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error query goals", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	filters := query.NewTransactionFilters(nil, []string{aid}, nil, "", "", "", "")

	transactions, err := h.transactionQuery.FindFilteredWithCount(
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
			views.NewAccountTransaction(t, bank.Name),
		)
	}

	var optionsLinks []components.BreadcrumbsLink
	for _, g := range banks {
		optionsLinks = append(optionsLinks, components.BreadcrumbsLink{
			Href: fmt.Sprintf("/banks/%s", g.ID),
			Name: g.Name,
		})
	}

	breadcrumbs := components.Breadcrumbs([]components.BreadcrumbsLink{
		{
			Href: "/banks",
			Name: "Banks",
		},
		{
			Href: "/banks/" + string(bank.ID),
			Name: bank.Name,
		},
	}, optionsLinks)

	content := templ.Join(
		breadcrumbs,
		partials.BankPage(*bank, transactionsViews, hasMore, transactions.Total),
	)

	renderWithMainLayout(w, r, bank.Name, content)
}
