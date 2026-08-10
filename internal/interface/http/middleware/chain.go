// Package middleware provides the global HTTP middleware chain (Recovery,
// RequestID, Logging, Auth) and the route-level guards (AuthProtectedRoute,
// AuthenticatedRedirect, RequireHTMX) used by the muxes wired in
// cmd/server/main.go.
package middleware

import (
	"net/http"
)

// Chain composes middlewares onto a handler, outermost-first: the first
// element of middlewares becomes the outermost wrapper.
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
