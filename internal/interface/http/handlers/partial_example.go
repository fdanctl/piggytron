package handlers

import (
	"net/http"
	"time"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/web/templates/components"
)

func ExampleHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		components.Button(
			time.Now().String(),
			"",
			components.BtnPrimary,
			components.BtnMedium,
			templ.Attributes{
				"hx-get":  "/partials/example",
				"hx-swap": "outerHTML",
			},
		).Render(r.Context(), w)

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func ExampleHandler2(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte(time.Now().String()))

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
