package middleware

import (
	"context"
	"fmt"
	"net/http"

	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
)

type contentKey string

const UserKey contentKey = "user"

func AuthMiddleware(store *rdb.SessionStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_id")
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			userInfo := store.Get(r.Context(), cookie.Value)
			if userInfo != nil {
				ctx := context.WithValue(r.Context(), UserKey, userInfo)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthProtectedRoute(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uinfo := r.Context().Value(UserKey)
		if uinfo == nil {
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
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}
