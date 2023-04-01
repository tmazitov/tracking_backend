package validation

import "errors"

var (
	ErrNotValid error = errors.New("value is not valid")
)
