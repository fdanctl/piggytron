package handlers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/web/templates/layouts"
)

var (
	ErrInvalidAmount = errors.New("invalid amount")
	ErrEmpty         = errors.New("string is empty")
)

const LIMIT = 30

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

func convertAmountStrToInt(str string) (int, error) {
	str = strings.ReplaceAll(str, ",", "")
	i := strings.Index(str, ".")
	tAmount := 0

	length := utf8.RuneCountInString(str)
	if str == "" {
		return 0, nil
	}

	if i == -1 {
		tAmount, err := strconv.Atoi(str)
		if err != nil {
			return 0, ErrInvalidAmount
		}
		tAmount *= 100
		return tAmount, nil
	}

	if length-1-i > 2 {
		return 0, ErrInvalidAmount
	}

	for length-i < 3 {
		str += "0"
		length++
	}

	tAmount, err := strconv.Atoi(strings.Replace(str, ".", "", 1))
	if err != nil {
		return 0, ErrInvalidAmount
	}

	return tAmount, nil
}
