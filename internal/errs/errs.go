// Package errs provides domain-specific error types and utilities.
package errs

import (
	"fmt"
	"net/http"
)

// ErrorType categorizes domain errors.
type ErrorType string

const (
	ErrNotFound     ErrorType = "not_found"
	ErrInvalidInput ErrorType = "invalid_input"
	ErrConflict     ErrorType = "conflict"
	ErrUnauthorized ErrorType = "unauthorized"
	ErrForbidden    ErrorType = "forbidden"
	ErrRateLimit    ErrorType = "rate_limit"
	ErrInternal     ErrorType = "internal_error"
)

// Error is a domain error with a type.
type Error interface {
	error
	Type() ErrorType
	Unwrap() error
}

// domainError implements Error.
type domainError struct {
	typ ErrorType
	msg string
	err error
}

func (e *domainError) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.msg, e.err)
	}
	return e.msg
}

func (e *domainError) Type() ErrorType { return e.typ }
func (e *domainError) Unwrap() error   { return e.err }

// New creates a new domain error.
func New(t ErrorType, msg string) Error {
	return &domainError{typ: t, msg: msg}
}

// Wrap wraps an existing error with a domain error type.
func Wrap(t ErrorType, msg string, err error) Error {
	return &domainError{typ: t, msg: msg, err: err}
}

// Is reports whether err is a domain error of the given type.
func Is(err error, t ErrorType) bool {
	if err == nil {
		return false
	}
	if de, ok := err.(Error); ok {
		return de.Type() == t
	}
	return false
}

// HTTPStatus maps error types to HTTP status codes.
func HTTPStatus(t ErrorType) int {
	switch t {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrInvalidInput:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	case ErrRateLimit:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// HTTPStatusFromError extracts HTTP status from a domain error.
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	if de, ok := err.(Error); ok {
		return HTTPStatus(de.Type())
	}
	return http.StatusInternalServerError
}
