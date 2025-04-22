package httpmapper

import (
	"net/http"

	"github.com/aghaghiamh/gocast/QAGame/pkg/richerr"
)

func MapResponseCustomErrorToHttp(err error) (code int, msg string) {
	switch err.(type) {
	case richerr.RichErr:
		re := err.(richerr.RichErr)
		statusCode := richErrCodeToHttpStatusCode(re.Code())

		if statusCode > 500 {
			// TODO: log the actual error
			statusCode = 500
		}

		return statusCode, re.Message()
	default:
		return http.StatusBadRequest, err.Error()
	}
}

func richErrCodeToHttpStatusCode(code int) int {
	switch code {
	case richerr.ErrServer:
		return http.StatusInternalServerError
	case richerr.ErrEntityNotFound:
		return http.StatusNotFound
	case richerr.ErrInvalidToken:
		return http.StatusBadRequest
	case richerr.ErrEntityDuplicate:
		return http.StatusConflict
	case richerr.ErrInvalidInput:
		return http.StatusUnprocessableEntity
	case richerr.ErrUnexpected:
		return http.StatusInternalServerError
	case richerr.ErrUnauthorized:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
