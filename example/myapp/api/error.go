package api

import (
	"github.com/pkg/errors"
)

type (
	Error struct {
		Err        error
		StatusCode int
		Message    string
	}
)

func (e *Error) Error() string {
	return e.Message
}

func NewError(err error, message string, status int) *Error {
	return &Error{
		Err:        err,
		Message:    message,
		StatusCode: status,
	}
}

func NewErrorWrap(err error, prefix, suffix, message string, status int) *Error {
	return &Error{
		Err:        errors.Wrapf(err, "%s/%s", prefix, suffix),
		Message:    message,
		StatusCode: status,
	}
}
