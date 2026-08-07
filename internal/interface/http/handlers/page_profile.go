package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

type ProfileHandler struct {
	userService *appuser.Service
}

func NewProfileHandler(us *appuser.Service) *ProfileHandler {
	return &ProfileHandler{
		userService: us,
	}
}

func (h *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProfileHandler) Get(w http.ResponseWriter, r *http.Request) {
	pf := views.NewProfileForm("User")
	cpf := views.NewChangePasswordForm()
	content := pages.Profile(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Account settings"},
			},
			Options: nil,
		},
		*pf,
		*cpf,
	)
	renderWithMainLayout(w, r, "Categories", content)
}
