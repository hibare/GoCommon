package http

import (
	"encoding/json"
	"net/http"

	"github.com/hibare/GoCommon/pkg/errors"
)

func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	e := errors.Error{
		Code:    statusCode,
		Message: err.Error(),
	}

	WriteJsonResponse(w, statusCode, e)
}
