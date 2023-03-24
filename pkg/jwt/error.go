package jwt

import "errors"

var (
	ErrInvalidToken             = errors.New("token is invalid")
	ErrInvalidTokenConvertation = errors.New("invalid token convertation")
	ErrInvalidClaims            = errors.New("token claims is invalid")
	ErrUnauthorized             = errors.New("user is not authorized")
	ErrTokenIsNotExist          = errors.New("token is not exist")
)
