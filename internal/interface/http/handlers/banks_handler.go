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
}

func NewBanksHandler(
	s *accountapp.Service,
	tq query.TransactionQueryService,
) *BanksHandler {
	return &BanksHandler{
		service:          s,
		transactionQuery: tq,
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

	banks, err := h.service.ReadAllBanksByUser(r.Context(), sessionInfo.UserID)
	if err != nil {
		logger.Error("error reading all banks", "error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	var bviews []views.Bank
	for _, v := range banks {
		bviews = append(bviews, views.NewBank(v))
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

	var tviews []views.Transaction
	for _, v := range transactions {
		tviews = append(tviews, views.NewTransaction(v))
	}

	content := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		err := components.Breadcrumbs([]components.BreadcrumbsLink{
			{Href: "", Name: "Banks"},
		}, nil).Render(ctx, w)
		if err != nil {
			return err
		}

		err = partials.Banks(bviews, tviews).Render(ctx, w)
		return err
	})

	renderWithMainLayout(w, r, "Banks", content)
}

func (h *BanksHandler) GetWithID(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	aid := r.PathValue("id")
	bank, err := h.service.ReadOneByID(r.Context(), aid)
	if err != nil {
		logger.Error("error read one bank", "error", err, "account_id", aid)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	fmt.Fprintf(w, "<p>%s</p>", bank.Name())
	fmt.Fprintf(w, "<p>%s</p>", bank.Currency())
}
