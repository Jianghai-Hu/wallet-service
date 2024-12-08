package utils

import "errors"

var _ error = (*myError)(nil)

type myError struct {
	code int
	msg  string
}

func (e *myError) Error() string {
	return e.msg
}

func (e *myError) Code() int {
	return e.code
}

func NewMyError(code int, msg string) *myError {
	return &myError{code: code, msg: msg}
}

func WrapMyError(code int, err error) *myError {
	if err == nil {
		return nil
	}
	var e *myError
	if errors.As(err, &e) {
		return e
	}
	return NewMyError(code, err.Error())
}
