package errs

import (
	"net/http"
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
	return NewError(http.StatusNotFound, msg, err)
}

func BadRequest(msg string, err error) *AppError {
	return NewError(http.StatusBadRequest, msg, err)
}

func Unauthorized(msg string, err error) *AppError {
	return NewError(http.StatusUnauthorized, msg, err)
}

func Internal(err error) *AppError {
	return NewError(http.StatusInternalServerError, "internal server error", err)
}
