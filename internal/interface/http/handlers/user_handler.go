package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	userapp "github.com/fdanctl/piggytron/internal/application/user"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type UserHandler struct {
	service *userapp.Service
}

func NewUserHandler(s *userapp.Service) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	action := r.PathValue("action")
	if action == "" {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	switch action {

	case "login":
		switch r.Method {

		case http.MethodPost:
			name := r.FormValue("name")
			pwd := r.FormValue("password")
			view := views.LoginView{
				Initial:  false,
				Name:     name,
				Password: pwd,
				Error:    "",
			}
			msgs := view.Validate()

			if len(msgs) == 0 {
				err := h.service.LoginUser(r.Context(), name, pwd)
				if err != nil {
					if errors.Is(err, sql.ErrNoRows) || errors.Is(err, userapp.ErrWrongPassword) {
						view.Error = "Name or password are invalid"
					}
				}
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
			partials.LoginForm(view).Render(r.Context(), w)

		default:
			http.NotFound(w, r)
		}

	case "signin":

	default:
		http.NotFound(w, r)

	}
}
