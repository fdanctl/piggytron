package middleware

import (
	"net/http"
	"runtime/debug"
)

// RecoveryMiddleware catches panics, logs the stack trace and answers with a
// generic 500 instead of crashing the process.
func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := LoggerFromContext(r.Context())

		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					"error", err,
					"stack", string(debug.Stack()),
				)

				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
