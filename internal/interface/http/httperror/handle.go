package httperror

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/a-h/templ"
	"github.com/fdanctl/piggytron/internal/errs"
	"github.com/fdanctl/piggytron/internal/interface/http/middleware"
)

func SendError(w http.ResponseWriter, r *http.Request, err error) bool {
	if err == nil {
		return false
	}

	logger := middleware.LoggerFromContext(r.Context())

	var apperr *errs.AppError
	if errors.As(err, &apperr) {
		if apperr.Kind == errs.KindInternal {
			logger.Error("request failed",
				"kind", apperr.Kind,
				"op", apperr.Operation,
				"message", apperr.Message,
				"error", apperr.Err,
			)
		} else {
			logger.Debug("bad request",
				"kind", apperr.Kind,
				"op", apperr.Operation,
				"message", apperr.Message,
				"error", apperr.Err,
			)
		}
		sendTrigger(w, apperr.Message)
		w.WriteHeader(HTTPStatus(apperr.Kind))
		return true
	}

	// Unknown error — generic
	logger.Error("unhandled error", "error", err)
	sendTrigger(w, "Something went wrong")
	w.WriteHeader(http.StatusInternalServerError)
	return true
}

func SendFormError(w http.ResponseWriter, r *http.Request, err error, form templ.Component) {
	logger := middleware.LoggerFromContext(r.Context())

	var apperr *errs.AppError
	if errors.As(err, &apperr) {
		if apperr.Kind == errs.KindInternal {
			logger.Error("unexpected error",
				"kind", apperr.Kind,
				"op", apperr.Operation,
				"message", apperr.Message,
				"error", apperr.Err,
			)
		} else {
			logger.Debug("bad request",
				"kind", apperr.Kind,
				"op", apperr.Operation,
				"message", apperr.Message,
				"error", apperr.Err,
			)
		}
	} else {
		logger.Error("unhandled error with form", "error", err)
		apperr = errs.NewInternalAppError(err, "httperr.SendFormError")
	}

	if apperr.Kind != errs.KindValidation {
		sendTrigger(w, apperr.Message)
	}
	w.WriteHeader(HTTPStatus(apperr.Kind))
	form.Render(r.Context(), w)
}

func sendTrigger(w http.ResponseWriter, msg string) {
	trigger := map[string]any{
		"show-toast": map[string]any{
			"level":   "error",
			"message": msg,
		},
	}
	b, _ := json.Marshal(trigger)
	w.Header().Set("HX-Trigger", string(b))
}
