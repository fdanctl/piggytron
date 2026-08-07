package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

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

func (h *AccountHandler) Get(w http.ResponseWriter, r *http.Request) {
	pf := views.NewProfileForm("User")
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
