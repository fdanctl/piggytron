package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

// AccountHandler renders the account settings page (profile and password
// forms).
type AccountHandler struct {
	userService *appuser.Service
}

func NewAccountHandler(us *appuser.Service) *AccountHandler {
	return &AccountHandler{
		userService: us,
	}
}

func (h *AccountHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the account settings page for the authenticated user.
func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	u, err := h.userService.FindByID(r.Context(), sessionInfo.UserID)
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	pf := views.NewProfileForm(u.Name())
	cpf := views.NewChangePasswordForm()
	content := pages.Account(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Account settings"},
			},
			Options: nil,
		},
		*pf,
		*cpf,
	)
	renderWithMainLayout(w, r, "Account", content)
}
