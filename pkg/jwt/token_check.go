package jwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

func (s *JwtStorage) ValidateCheck(ctx context.Context, token string) (*CheckClaims, error) {

	// Check if token is exists
	if err := s.CheckIsExists(ctx, token); err != nil {
		return nil, err
	}

	return s.verifyCheckToken(ctx, token)
}

func (s *JwtStorage) CheckIsExists(ctx context.Context, token string) error {
	return s.isExists(ctx, "che:", token)
}

func (s *JwtStorage) verifyCheckToken(ctx context.Context, token string) (*CheckClaims, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return s.config.Secret, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &CheckClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*CheckClaims)
	if !ok {
		return nil, ErrInvalidTokenConvertation
	}

	return payload, nil
}
