package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

type ctxKey string

const UserKey ctxKey = "user"

func AuthMiddleware(store *rdb.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := LoggerFromContext(r.Context())
			cookie, err := r.Cookie("session_id")
			if err != nil {
				logger.Info("unauthenticated - no session cookie")
				next.ServeHTTP(w, r)
				return
			}

			userInfo := store.Get(r.Context(), cookie.Value)
			if userInfo != nil {
				ctx := context.WithValue(r.Context(), UserKey, userInfo)

				logger = logger.With(
					slog.String("user_id", userInfo.UserID),
				)

				ctx = context.WithValue(ctx, LoggerKey, logger)

				logger.Info("authenticated")
				r = r.WithContext(ctx)

			} else {
				logger.Info("unauthenticated - expired")
			}

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
