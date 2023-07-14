package errors

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
