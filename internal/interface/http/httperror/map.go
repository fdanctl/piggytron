package httperror

import (
	"net/http"

	"github.com/fdanctl/piggytron/internal/errs"
)

func HTTPStatus(k errs.ErrorKind) int {
	switch k {
	case errs.KindValidation:
		return http.StatusUnprocessableEntity // 422

	case errs.KindBusinessRule:
		return http.StatusUnprocessableEntity // 422

	case errs.KindNotFound:
		return http.StatusNotFound // 404

	case errs.KindConflict:
		return http.StatusConflict // 409

	case errs.KindBadRequest:
		return http.StatusBadRequest // 400

	case errs.KindUnauthorized:
		return http.StatusUnauthorized // 401

	case errs.KindInternal:
		return http.StatusInternalServerError // 500

	default:
		return http.StatusInternalServerError
	}
}
