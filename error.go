package pam4sdk

import (
	"runtime"
)

// IError is interface for error
type IError interface {
	Error() string
}

// Error is the struct for error
type Error struct {
	message string
	err     error
	forUser bool
}

// NewErr return new error without logging
func NewErr(err error) *Error {
	e, ok := err.(*Error)
	isDoubleWrapSameError := ok && e.err != nil && e.err.Error() == err.Error()
	if isDoubleWrapSameError {
		return e
	}
	return &Error{
		message: err.Error(),
		err:     err,
	}
}

// NewErrM return error message without log
func NewErrM(message string) *Error {
	return &Error{
		message: message,
		err:     nil,
	}
}

// NewErrorE return error and log error message to logger
func NewErrorE(logger ILogger, err error) *Error {
	e, ok := err.(*Error)
	isDoubleWrapSameError := ok && e.err != nil && e.err.Error() == err.Error()
	if isDoubleWrapSameError {
		return e
	}
	// Log error then wrap error into Error
	_, fn, line, _ := runtime.Caller(1)
	logger.ErrorFL(err.Error(), fn, line)
	return &Error{
		message: err.Error(),
		err:     err,
	}
}

// NewErrorM return error and log message to logger
func NewErrorM(logger ILogger, message string) *Error {
	_, fn, line, _ := runtime.Caller(1)
	logger.ErrorFL(message, fn, line)
	return &Error{
		message: message,
		err:     nil,
	}
}

// NewErrorEU return error for user and log error message to logger
func NewErrorEU(logger ILogger, err error) *Error {
	e := NewErrorE(logger, err)
	e.forUser = true
	return e
}

// NewErrorMU return error for user and log message to logger
func NewErrorMU(logger ILogger, message string) *Error {
	e := NewErrorM(logger, message)
	e.forUser = true
	return e
}

// Error return error message
func (e *Error) Error() string {
	return e.message
}
