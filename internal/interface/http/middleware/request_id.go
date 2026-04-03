package middleware

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const RequestIDKey ctxKey = "request_id"

func generateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 36) + "-" +
		strconv.Itoa(rand.Intn(100000))
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get("X-Request-ID")
		if reqID == "" {
			reqID = generateRequestID()
		}
		w.Header().Set("X-Request-ID", reqID)
		ctx := context.WithValue(r.Context(), RequestIDKey, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
