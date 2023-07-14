package errors

import "errors"

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
