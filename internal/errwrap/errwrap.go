package errwrap

import (
	"errors"
	"fmt"
)

// Wrap is packing all errors to one
// also (all) nil errors are handled
func Wrap(errors ...error) error {
	var err error

	for _, nextErr := range errors {
		if nextErr == nil {
			continue
		}
		if err == nil {
			// first error
			err = nextErr
			continue
		}
		// next errors
		err = fmt.Errorf("%w; %s", err, nextErr.Error())
	}

	return err
}

// FirstError unwraps the first error
func FirstError(err error) error {
	for {
		unwrappedErr := errors.Unwrap(err)
		if unwrappedErr == nil {
			return err
		}
		err = unwrappedErr
	}
}
