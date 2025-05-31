// Package http provides HTTP utilities for the application.
package http

import (
	"encoding/json"
	"net/http"

	"github.com/hibare/GoCommon/v2/pkg/errors"
)

// WriteJsonResponse writes a JSON response with the given status code and data.
func WriteJsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

// WriteErrorResponse writes an error response with the given status code and error.
func WriteErrorResponse(w http.ResponseWriter, statusCode int, err error) {
	e := errors.Error{
		Code:    statusCode,
		Message: err.Error(),
	}

	WriteJsonResponse(w, statusCode, e)
}
