package middleware

import "net/http"

func RequireHTMX(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			http.NotFound(w, r)
			return
		}
		next(w, r)
	}
}
