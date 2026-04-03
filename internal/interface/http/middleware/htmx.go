package middleware

import "net/http"

func RequireHTMX(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("HX-Request") != "true" {
			logger := LoggerFromContext(r.Context())
			logger.Warn("non-HTMX request in a HTMX only route")
			http.NotFound(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}
