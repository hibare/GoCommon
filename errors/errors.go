package errors

import "errors"

var (
	ErrUnauthorized = errors.New("unauthorized")
)

type Error struct {
	Message string `json:"message"`
}
