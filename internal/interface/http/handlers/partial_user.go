package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/domain/user"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/interface/http/shared"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type UserHandler struct {
	service     *appuser.Service
	cookieMaker *shared.CookieMaker
}

func NewUserHandler(s *appuser.Service, cm *shared.CookieMaker) *UserHandler {
	return &UserHandler{
		service:     s,
		cookieMaker: cm,
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
	logger := middleware.LoggerFromContext(r.Context())
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
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.LoginForm(view).Render(r.Context(), w)
		return
	}

	sid, err := h.service.LoginUser(r.Context(), name, pwd)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, appuser.ErrWrongPassword) {
			logger.Info("error on login", "error", err)
			view.ErrorMsg = "Name or password are invalid"
			w.WriteHeader(http.StatusUnprocessableEntity)
			partials.LoginForm(view).Render(r.Context(), w)
			return
		}
		logger.Error("error on login", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(redirect)
	if err != nil || u.IsAbs() || u.Host != "" {
		redirect = "/"
	}

	http.SetCookie(w, h.cookieMaker.NewCookie(sid))
	w.Header().Set("HX-Redirect", redirect)
	w.WriteHeader(http.StatusSeeOther)
}

func (h *UserHandler) SignupPost(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
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
		logger.Info("invalid form", "error", msgs)
		w.WriteHeader(http.StatusUnprocessableEntity)
		partials.SignupForm(view).Render(r.Context(), w)
		return
	}

	sid, err := h.service.CreateUser(r.Context(), name, pwd)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		if errors.Is(err, user.ErrDuplicate) {
			logger.Info("error on signup", "error", err)
			view.CustomError = err
		} else {
			logger.Error("error on signup", "error", err)
		}
		partials.SignupForm(view).Render(r.Context(), w)
		return
	}

	http.SetCookie(w, h.cookieMaker.NewCookie(sid))
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) LogoutGet(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	cookie, _ := r.Cookie("session_id")

	err := h.service.LogoutUser(r.Context(), cookie.Value)
	if err != nil {
		logger.Error("error on logout", "error", err)
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, h.cookieMaker.RevokeCookie())
	w.Header().Set("HX-Redirect", "")
	w.WriteHeader(http.StatusNoContent)
}
