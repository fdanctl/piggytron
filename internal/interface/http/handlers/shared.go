package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/a-h/templ"
	rdb "github.com/fdanctl/piggytron/internal/infrastructure/redis"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

const LIMIT = 30

var ErrInvalidSession = errors.New("invalid session")

func sessionInfoFromCtx(ctx context.Context) (*rdb.SessionInfo, error) {
	v := ctx.Value(middleware.UserKey)
	if v == nil {
		return nil, ErrInvalidSession
	}

	sessionInfo, ok := v.(*rdb.SessionInfo)
	if !ok {
		return nil, ErrInvalidSession
	}
	return sessionInfo, nil
}

func renderWithMainLayout(
	w http.ResponseWriter,
	r *http.Request,
	title string,
	content templ.Component,
) error {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("Hx-Request") == "true" {
		err := content.Render(r.Context(), w)
		fmt.Fprintf(w, "<title>%s</title>", title)
		return err
	}

	main := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		ctx = templ.WithChildren(ctx, content)
		err := layouts.Main().Render(ctx, w)
		return err
	})

	ctx := templ.WithChildren(r.Context(), main)
	return layouts.Base(title).Render(ctx, w)
}
