package httputils

import (
	"net/http"
)

type Error struct {
	StatusCode int
	Err        error
	Message    string
}

func NewBadRequestError(err error, msg string) Error {
	return Error{
		StatusCode: http.StatusBadRequest,
		Err:        err,
		Message:    msg,
	}
}

func NewInternalServerError(err error, msg string) Error {
	return Error{
		StatusCode: http.StatusInternalServerError,
		Err:        err,
		Message:    msg,
	}
}

func NewNotFoundError(err error, msg string) Error {
	return Error{
		StatusCode: http.StatusNotFound,
		Err:        err,
		Message:    msg,
	}
}
