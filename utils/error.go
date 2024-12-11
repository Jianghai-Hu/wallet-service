package utils

import (
	"errors"

	"jianghai-hu/wallet-service/internal/common"
)

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

func NewMyError(code int, msg string) error {
	return &myError{code: code, msg: msg}
}

func WrapMyError(code int, err error) error {
	if err == nil {
		return nil
	}

	var e *myError
	if errors.As(err, &e) {
		return e
	}

	return NewMyError(code, err.Error())
}

func ResolveError(err error) (int, string) {
	if err == nil {
		return 0, ""
	}

	var e *myError
	if errors.As(err, &e) {
		return e.Code(), e.Error()
	}

	return common.Constant_ERROR_UNKNOW, err.Error()
}
