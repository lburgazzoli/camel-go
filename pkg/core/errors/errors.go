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
	return NotImplemented(fmt.Sprintf(format, args))
}

func IsNotImplemented(err error) bool {
	return errors.Is(err, &notImplemented{})
}
