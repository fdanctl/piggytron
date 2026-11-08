package handlers

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/interface/http/assets"
	"github.com/fdanctl/piggytron/web/templates/pages"
	"github.com/fdanctl/piggytron/web/views"
)

// PreferencesHandler renders the preferences page (region/formatting
// settings).
type PreferencesHandler struct {
	userService *appuser.Service
	assetsURLs  assets.URLs
}

func NewPreferencesHandler(us *appuser.Service, assetsURLs assets.URLs) *PreferencesHandler {
	return &PreferencesHandler{
		userService: us,
		assetsURLs:  assetsURLs,
	}
}

func (h *PreferencesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.Get(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Get renders the preferences page.
func (h *PreferencesHandler) Get(w http.ResponseWriter, r *http.Request) {
	rf := views.NewRegionForm()
	content := pages.Preferences(
		views.BreadcrumbsView{
			Items: []views.BreadcrumbsLink{
				{Href: "", Name: "Account settings"},
			},
			Options: nil,
		},
		*rf,
	)
	renderWithMainLayout(w, r, "Preferences", h.assetsURLs, content)
}
