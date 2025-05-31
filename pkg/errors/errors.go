// Package errors provides common error variables for the application.
package errors

import "errors"

var (
	// ErrUnauthorized indicates an unauthorized action.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInternalServerError indicates an internal server error.
	ErrInternalServerError = errors.New("internal server error")

	// ErrNonOKError indicates a non-OK error.
	ErrNonOKError = errors.New("non ok error")

	// ErrChecksumMismatch indicates a checksum mismatch.
	ErrChecksumMismatch = errors.New("checksum Mismatch")

	// ErrCreatingDir indicates a directory creation error.
	ErrCreatingDir = errors.New("creating dir")

	// ErrNotDir indicates a not a directory error.
	ErrNotDir = errors.New("not a directory")

	// ErrNotFile indicates a not a file error.
	ErrNotFile = errors.New("not a file")

	// ErrRecordNotFound indicates a record not found error.
	ErrRecordNotFound = errors.New("record not found")
)

// Error represents an error with a code and message.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
