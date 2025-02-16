package errors

import "errors"

var (
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInternalServerError = errors.New("internal server error")
	ErrNonOKError          = errors.New("non ok error")
	ErrChecksumMismatch    = errors.New("checksum Mismatch")
	ErrCreatingDir         = errors.New("creating dir")
	ErrNotDir              = errors.New("not a directory")
	ErrNotFile             = errors.New("not a file")
	ErrRecordNotFound      = errors.New("record not found")
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
