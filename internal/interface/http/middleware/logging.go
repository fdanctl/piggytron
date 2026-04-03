package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

const LoggerKey ctxKey = "logger"

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(l *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK, // default
			}

			logger := l.With(
				slog.String("request_id", RequestIDFromContext(r.Context())),
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
			logger.Info("request start")

			ctx := context.WithValue(r.Context(), LoggerKey, logger)

			next.ServeHTTP(rw, r.WithContext(ctx))

			logger = LoggerFromContext(ctx)

			if rw.status >= 500 {
				logger.Error("request failed",
					slog.Int("status", rw.status),
					slog.Duration("duration", time.Since(start)),
				)
			} else {
				logger.Info("request completed",
					slog.Int("status", rw.status),
					slog.Duration("duration", time.Since(start)),
				)
			}
		})
	}
}
