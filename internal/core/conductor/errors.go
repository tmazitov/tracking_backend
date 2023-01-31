package conductor

import "errors"

var ErrTooManyAttempts error = errors.New("too many attempts to auth")
