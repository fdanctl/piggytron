package handlers

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/layouts"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

// UserHandler handles the profile form submission (name change).
type UserHandler struct {
	service *appuser.Service
}

func NewUserHandler(s *appuser.Service) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Post(w, r)

	default:
		http.NotFound(w, r)
	}
}

// Post validates and applies the profile form, re-rendering it with a
// success toast on completion.
func (h *UserHandler) Post(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.ProfileForm{
		Name: r.FormValue("name"),
	}

	view.Initial = false
	msgs := view.Validate()
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.ProfileForm(view).Render(r.Context(), w)
		return
	}

	err = h.service.ChangeName(r.Context(), sessionInfo.UserID, cookie.Value, view.Name)
	if err != nil {
		view.SetError(err)
		httperror.SendFormError(w, r, err, partials.ProfileForm(view))
		return
	}

	username := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err = fmt.Fprint(w, view.Name)
		return err
	})

	templ.Join(
		partials.ProfileForm(view),
		layouts.HxPartial("#sidebar-username", "innerHTML", username),
		components.SendToast(
			components.Success,
			"User name updated",
		),
	).Render(r.Context(), w)
}
