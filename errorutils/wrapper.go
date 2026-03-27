package errorutils

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

// WrapIfNotNil the lazy persons wrapper, if there is a nil for err, it returns nil,
// if error has a value, then the error is wrapped with the message and returned
func WrapIfNotNil(message string, err error) error {
	return WrapError(err, message)
}

// ContainsErrorSubstring checks if the error or any of its wrapped errors contain the target substring.
func ContainsErrorSubstring(err error, target string) bool {
	for err != nil {
		if strings.Contains(err.Error(), target) {
			return true
		}
		err = errors.Unwrap(err)
	}
	return false
}

func WrapError(err error, context ...string) error {
	if err == nil {
		return nil
	}

	callerName := "unknown"
	if pc, _, _, ok := runtime.Caller(1); ok {
		if fn := runtime.FuncForPC(pc); fn != nil {
			callerName = fn.Name()
		}
	}

	parts := make([]string, 0, 1+len(context))
	parts = append(parts, callerName)
	parts = append(parts, context...)

	return fmt.Errorf("%s: %w", strings.Join(parts, " - "), err)
}
