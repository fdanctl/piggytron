package handlers

import (
	"net/http"
	"net/url"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/application/appuser"
	"github.com/fdanctl/piggytron/internal/interface/http/httperror"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/internal/interface/http/shared"
	"github.com/fdanctl/piggytron/web/templates/components"
	"github.com/fdanctl/piggytron/web/templates/partials"
	"github.com/fdanctl/piggytron/web/views"
)

type AuthHandler struct {
	service     *appuser.Service
	cookieMaker *shared.CookieMaker
}

func NewAuthHandler(s *appuser.Service, cm *shared.CookieMaker) *AuthHandler {
	return &AuthHandler{
		service:     s,
		cookieMaker: cm,
	}
}

func (h *AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	case "change-password":
		switch r.Method {
		case http.MethodPost:
			h.PostPassword(w, r)

		default:
			http.NotFound(w, r)
		}

	default:
		http.NotFound(w, r)

	}
}

func (h *AuthHandler) LoginPost(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	redirect := r.FormValue("redirect")
	view := views.LoginView{
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
		view.SetError(err)
		form := partials.LoginForm(view)
		httperror.SendFormError(w, r, err, form)
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

func (h *AuthHandler) SignupPost(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	name := r.FormValue("name")
	pwd := r.FormValue("password")
	pwdConf := r.FormValue("password-confirm")
	view := views.SignupView{
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
		view.SetError(err)
		form := partials.SignupForm(view)
		httperror.SendFormError(w, r, err, form)
		return
	}

	http.SetCookie(w, h.cookieMaker.NewCookie(sid))
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) LogoutGet(w http.ResponseWriter, r *http.Request) {
	logger := middleware.LoggerFromContext(r.Context())
	cookie, _ := r.Cookie("session_id")

	if r.URL.Query().Get("all") == "" {
		err := h.service.LogoutUser(r.Context(), cookie.Value)
		if err != nil {
			logger.Error("error on logout", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	} else {
		sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
		if err != nil {
			httperror.SendError(w, r, err)
			return
		}

		err = h.service.LogoutUserFromAllDevices(r.Context(), sessionInfo.UserID, cookie.Value)
		if err != nil {
			logger.Error("error on logout from all devices", "error", err)
			http.Error(w, "", http.StatusInternalServerError)
			return
		}
	}

	http.SetCookie(w, h.cookieMaker.RevokeCookie())
	w.Header().Set("HX-Redirect", "")
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) PostPassword(w http.ResponseWriter, r *http.Request) {
	sessionInfo, err := middleware.SessionInfoFromCtx(r.Context())
	if err != nil {
		httperror.SendError(w, r, err)
		return
	}

	view := views.ChangePasswordForm{
		CurrentPassword:    r.FormValue("current-password"),
		NewPassword:        r.FormValue("new-password"),
		NewPasswordConfirm: r.FormValue("new-password-confirm"),
	}
	msgs := view.Validate()
	form := partials.ChangePasswordForm(view)

	if len(msgs) > 0 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		form.Render(r.Context(), w)
		return
	}

	sid, err := h.service.ChangePassword(
		r.Context(),
		sessionInfo.UserID,
		view.CurrentPassword,
		view.NewPassword,
	)
	if err != nil {
		view.SetError(err)
		form = partials.ChangePasswordForm(view)
		httperror.SendFormError(w, r, err, form)
		return
	}

	http.SetCookie(w, h.cookieMaker.NewCookie(sid))
	templ.Join(
		partials.ChangePasswordForm(*views.NewChangePasswordForm()),
		components.SendToast(
			components.Success,
			"Password updated",
		),
	).Render(r.Context(), w)
}
