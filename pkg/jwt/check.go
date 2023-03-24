package jwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
)

func (s *JwtStorage) verifyToken(ctx context.Context, token string) (*AccessClaims, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}

		return s.config.Secret, nil
	}

	jwtToken, err := jwt.ParseWithClaims(token, &AccessClaims{}, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrInvalidToken) {
			return nil, ErrInvalidToken
		}
		return nil, ErrInvalidToken
	}

	payload, ok := jwtToken.Claims.(*AccessClaims)
	if !ok {
		return nil, ErrInvalidTokenConvertation
	}

	return payload, nil
}

func (s *JwtStorage) isExists(ctx context.Context, prefix string, token string) error {

	err := s.redis.Get(ctx, prefix+token).Err()
	if err == redis.Nil {
		return ErrTokenIsNotExist
	}

	if err != nil {
		return err
	}

	return nil
}
