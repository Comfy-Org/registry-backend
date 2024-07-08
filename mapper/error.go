package mapper

import "errors"

type errCode uint8

const (
	errCodeUnknown errCode = iota
	errCodeBadRequest
)

type codedError struct {
	code errCode
	msg  string
	err  []error
}

func (e codedError) Error() string {
	return e.msg
}

func (e codedError) Unwrap() []error {
	return e.err
}

func isCodedError(err error, code errCode) bool {
	var e codedError
	if !errors.As(err, &e) {
		return false
	}

	return e.code == code
}

func NewErrorBadRequest(msg string, err ...error) error {
	return codedError{code: errCodeBadRequest, msg: msg, err: err}
}

func IsErrorBadRequest(err error) (y bool) {
	return isCodedError(err, errCodeBadRequest)
}
