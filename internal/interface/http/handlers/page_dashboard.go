package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/query"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type DashboardHandler struct {
	ledger query.LedgerQueryService
}

func NewDashboardHandler(lq query.LedgerQueryService) *DashboardHandler {
	return &DashboardHandler{
		ledger: lq,
	}
}

func (h *DashboardHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	transactions, err := h.ledger.GetRecentEntries(
		r.Context(),
		sessionInfo.UserID,
		5,
	)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	var tviews []views.Transaction
	for _, v := range transactions {
		tviews = append(tviews, views.NewTransaction(v))
	}

	content := pages.Dashboard(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Dashboard"},
			},
			Options: nil,
		},
		tviews,
	)

	renderWithMainLayout(w, r, "Dashboard", content)
}
