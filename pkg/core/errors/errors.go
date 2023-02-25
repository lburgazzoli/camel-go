package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

type commonError struct {
	message string
}

func (e *commonError) Error() string {
	return e.message
}

//
// Not Implemented
//

type notImplemented struct {
	commonError
}

func NotImplemented(message string) error {
	return &notImplemented{
		commonError: commonError{
			message: message,
		},
	}
}

func NotImplementedf(format string, args ...any) error {
	return NotImplemented(fmt.Sprintf(format, args...))
}

func IsNotImplemented(err error) bool {
	return errors.Is(err, &notImplemented{})
}

//
// Missing parameter
//

type missingParameter struct {
	commonError

	parameter string
}

func MissingParameter(parameter string, message string) error {
	answer := &missingParameter{
		parameter: parameter,
	}

	if message == "" {
		answer.message = fmt.Sprintf("mandatory parameter %s not set", parameter)
	} else {
		answer.message = fmt.Sprintf("%s: mandatory parameter %s not set", answer, parameter)
	}

	return answer
}

func MissingParameterf(parameter string, format string, args ...any) error {
	return MissingParameter(parameter, fmt.Sprintf(format, args...))
}

func IsMissingParameter(err error) bool {
	return errors.Is(err, &missingParameter{})
}

//
// Missing parameter
//

type invalidParameter struct {
	commonError

	parameter string
}

func InvalidParameter(parameter string, message string) error {
	answer := &invalidParameter{
		parameter: parameter,
	}

	if message == "" {
		answer.message = fmt.Sprintf("invalid parameter %s", parameter)
	} else {
		answer.message = fmt.Sprintf("%s: invalid parameter %s", answer, parameter)
	}

	return answer
}

func InvalidParameterf(parameter string, format string, args ...any) error {
	return InvalidParameter(parameter, fmt.Sprintf(format, args...))
}

func IsInvalidParameter(err error) bool {
	return errors.Is(err, &invalidParameter{})
}

//
// Internal error
//

type internalError struct {
	commonError
}

func InternalError(message string) error {
	answer := &internalError{
		commonError: commonError{
			message: message,
		},
	}

	return answer
}

func InternalErrorf(format string, args ...any) error {
	return InternalError(fmt.Sprintf(format, args...))
}

func IsInternalError(err error) bool {
	return errors.Is(err, &internalError{})
}

//
// Not found
//

type notFound struct {
	commonError
}

func NotFound(message string) error {
	answer := &notFound{
		commonError: commonError{
			message: message,
		},
	}

	return answer
}

func NotFoundf(format string, args ...any) error {
	return NotFound(fmt.Sprintf(format, args...))
}

func IsNotFound(err error) bool {
	return errors.Is(err, &notFound{})
}
