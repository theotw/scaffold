package errorutils

import (
	"fmt"
)

// WrapIfNotNil the lazy persons wrapper, if there is a nil for err, it returns nil,
// if error has a value, then the error is wrapped with the message and returned
func WrapIfNotNil(message string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s - %w", message, err)
}
