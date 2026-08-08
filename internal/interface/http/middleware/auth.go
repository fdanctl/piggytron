package middleware

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/fdanctl/piggytron/internal/auth"
)

type ctxKey string

const UserKey ctxKey = "user"

func AuthMiddleware(sessionManager *auth.SessionManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := LoggerFromContext(r.Context())
			cookie, err := r.Cookie("session_id")
			if err != nil {
				logger.Info("unauthenticated - no session cookie")
				next.ServeHTTP(w, r)
				return
			}

			userInfo, err := sessionManager.GetSession(r.Context(), cookie.Value)
			if err != nil {
				if errors.Is(err, auth.ErrNotFound) {
					logger.Info("unauthenticated - expired")
				} else {
					logger.Error("redis failed to get session", "err", err)
				}
				next.ServeHTTP(w, r)
				return
			}

			valid, err := sessionManager.ValidateSession(r.Context(), userInfo)
			if err != nil {
				logger.Error("redis failed to get session version", "err", err)
				next.ServeHTTP(w, r)
				return
			}

			if !valid {
				logger.Info("unauthenticated - revoked")
				next.ServeHTTP(w, r)
				return
			}
			ctx := context.WithValue(r.Context(), UserKey, userInfo)

			logger = logger.With(
				slog.String("user_id", userInfo.UserID),
			)

			ctx = context.WithValue(ctx, LoggerKey, logger)

			logger.Info("authenticated")
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

func AuthProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uinfo := r.Context().Value(UserKey)
		if uinfo == nil {
			logger := LoggerFromContext(r.Context())
			logger.Info("redirect")
			http.Redirect(w, r, fmt.Sprint("/login?redirect=", r.RequestURI), http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AuthenticatedRedirect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uinfo := r.Context().Value(UserKey)
		if uinfo != nil {
			logger := LoggerFromContext(r.Context())
			logger.Info("redirect")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
