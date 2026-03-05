package middleware

import "net/http"

func RequireHTMX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
