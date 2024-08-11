package errorutils

import "errors"

var TimeoutError = errors.New("timeout")

var DupRecordError = errors.New("dup record")
