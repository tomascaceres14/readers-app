package errs

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrBadRequest    = errors.New("bad request")
	ErrUnauthorized  = errors.New("unauthorized")
	ErrInternal      = errors.New("internal")
	ErrAlreadyExists = errors.New("already exists")
	ErrBadURL        = errors.New("url does not contain ")
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e AppError) Error() string {
	return e.Message
}

func NewError(code int, msg string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func NotFound(msg string, err error) *AppError {
	return NewError(http.StatusNotFound, msg, fmt.Errorf("%w: %w", ErrNotFound, err))
}

func BadRequest(msg string, err error) *AppError {
	return NewError(http.StatusBadRequest, msg, fmt.Errorf("%w: %w", ErrBadRequest, err))
}

func Unauthorized(msg string, err error) *AppError {
	return NewError(http.StatusUnauthorized, msg, fmt.Errorf("%w: %w", ErrUnauthorized, err))
}

func Internal(err error) *AppError {
	return NewError(http.StatusInternalServerError, "internal server error", fmt.Errorf("%w: %w", ErrInternal, err))
}
