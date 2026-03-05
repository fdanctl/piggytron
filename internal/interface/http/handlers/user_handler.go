package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"time"

	userapp "github.com/fdanctl/piggytron/internal/application/user"
	"github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type UserHandler struct {
	service      *userapp.Service
	sessionStore *redis.SessionStore
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
		switch r.Method {

		case http.MethodGet:
			h.LogoutGet(w, r)

		default:
			http.NotFound(w, r)
		}

	default:
		http.NotFound(w, r)

	}
}

func (h *UserHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	redirect := r.FormValue("redirect")
	view := views.LoginView{
		Initial:  false,
		Redirect: redirect,
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

	sid, err := h.service.LoginUser(r.Context(), name, pwd)
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

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})

	if redirect == "" || redirect[0] != '/' {
		redirect = "/"
	}
	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(http.StatusSeeOther)
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

	sid, err := h.service.CreateUser(r.Context(), name, pwd)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, userapp.ErrUserExists) {
			view.CustomError = err
			partials.SignupForm(view).Render(r.Context(), w)
			return
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sid,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Now().Add(time.Hour * 24),
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) LogoutGet(w http.ResponseWriter, r *http.Request) {
	err := h.service.LogoutUser(r.Context())
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Path:    "/",
		Expires: time.Now(),
	})
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}
