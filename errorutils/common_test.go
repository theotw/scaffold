package errorutils

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	err := TimeoutError
	assert.True(t, errors.Is(err, TimeoutError))

	wrapped := fmt.Errorf("had a timeout doing something %w", TimeoutError)
	assert.True(t, errors.Is(wrapped, TimeoutError))

	wrapped = fmt.Errorf("index issue %w", DupRecordError)
	assert.True(t, errors.Is(wrapped, DupRecordError))

}
