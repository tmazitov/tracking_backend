package conductor

import "errors"

var ErrTooManyAttempts error = errors.New("too many attempts to auth")
var ErrInvalidToken error = errors.New("invalid check token to auth")
