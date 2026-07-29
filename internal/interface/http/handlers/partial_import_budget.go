package handlers

import (
	"fmt"
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appbudget"
	"github.com/fdanctl/piggytron/internal/domain/budget"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/components"
)

type ImportBudgetHandler struct {
	service *appbudget.Service
}

func NewImportBudgetHandler(s *appbudget.Service) *ImportBudgetHandler {
	return &ImportBudgetHandler{service: s}
}

func (h *ImportBudgetHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Post(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ImportBudgetHandler) Post(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}
	month := r.URL.Query().Get("month")
	bm, err := budget.ParseMonth(month)
	if err != nil {
		err := errs.NewAppError(
			errs.KindBadRequest,
			fmt.Sprintf("%s is not a valid month", month),
			fmt.Errorf("failed to parse month '%s': %w", month, err),
			"ImportBudgetHandler.Post",
		)
		httperror.SendError(w, r, err)
		return
	}

	count, err := h.service.CopyFromLastMonth(r.Context(), sessionInfo.UserID, bm)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	if count > 0 {
		w.Header().Set("HX-Trigger", "budget-import-completed")
		components.SendToast(
			components.Success,
			fmt.Sprintf("Imported  %d categories with success", count),
		).Render(r.Context(), w)
	} else {
		components.SendToast(
			components.Warning,
			"No budget last month",
		).Render(r.Context(), w)
	}
}
