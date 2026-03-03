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
			h.LoginPost(w, r)

		default:
			http.NotFound(w, r)
		}

	case "signup":
		switch r.Method {

		case http.MethodPost:
			h.SignupPost(w, r)

		default:
			http.NotFound(w, r)
		}

	case "logout":

	default:
		http.NotFound(w, r)

	}
}

func (h *UserHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	view := views.LoginView{
		Initial:  false,
		Name:     name,
		Password: pwd,
	}
	msgs := view.Validate()

	// invalid form
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.LoginForm(view).Render(r.Context(), w)
		return
	}

	err := h.service.LoginUser(r.Context(), name, pwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, userapp.ErrWrongPassword) {
			view.ErrorMsg = "Name or password are invalid"
			w.WriteHeader(http.StatusUnprocessableEntity)
			partials.LoginForm(view).Render(r.Context(), w)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// TODO set cookies
	w.Header().Set("HX-Redirect", "/app")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) SignupPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	pwdConf := r.FormValue("password-confirm")
	view := views.SignupView{
		Initial:         false,
		Name:            name,
		Password:        pwd,
		PasswordConfirm: pwdConf,
	}
	msgs := view.Validate()

	// invalid form
	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.SignupForm(view).Render(r.Context(), w)
		return
	}

	err := h.service.CreateUser(r.Context(), name, pwd)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, userapp.ErrUserExists) {
			view.CustomError = err
			partials.SignupForm(view).Render(r.Context(), w)
			return
		}
		return
	}

	// TODO set cookies
	w.Header().Set("HX-Redirect", "/app")
	w.WriteHeader(http.StatusNoContent)
}
