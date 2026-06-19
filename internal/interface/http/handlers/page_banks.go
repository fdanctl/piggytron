package handlers

import (
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appaccount"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type BanksHandler struct {
	service          *appaccount.Service
	transactionQuery query.TransactionQueryService
	accountQuery     query.AccountQueryService
}

func NewBanksHandler(
	s *appaccount.Service,
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
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	transactions, err := h.transactionQuery.GetRecentTransactions(
		r.Context(),
		sessionInfo.UserID,
		5,
	)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find transactions: %w", err))
		return
	}

	accounts, err := h.accountQuery.FindAllWithSum(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, fmt.Errorf("failed to find accounts: %w", err))
		return
	}

	pageView := views.NewBankPage(accounts, transactions)

	content := pages.Banks(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Banks"},
			},
			Options: nil,
		},
		pageView,
	)

	renderWithMainLayout(w, r, "Banks", content)
}

func (h *BanksHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	aid := r.PathValue("id")
	bank, err := h.accountQuery.FindWithSum(r.Context(), aid)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	banks, err := h.accountQuery.FindBanksIDNames(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
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
		httperror.SendError(w, r, fmt.Errorf("failed to find filtered transactions: %w", err))
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

	var optionsLinks []views.BreadcrumbsLink
	for _, g := range banks {
		optionsLinks = append(optionsLinks, views.BreadcrumbsLink{
			Href: fmt.Sprintf("/banks/%s", g.ID),
			Name: g.Name,
		})
	}

	content := pages.Bank(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Transactions"},
				{Href: "", Name: "All"},
			},
			Options: optionsLinks,
		}, *bank, transactionsViews, hasMore, transactions.Total,
	)

	renderWithMainLayout(w, r, bank.Name, content)
}
