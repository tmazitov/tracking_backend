package jwt

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type AccessClaims struct {
	UserId    int64     `json:"user_id"`
	RoleId    int       `json:"role_id"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (c *AccessClaims) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return ErrInvalidToken
	}
	return nil
}

type TokenPair struct {
	Access  string
	Refresh string
}

func (s *JwtStorage) createToken(claims AccessClaims, duration time.Duration) (string, error) {

	now := time.Now()

	claims.IssuedAt = now
	claims.ExpiredAt = now.Add(duration)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	return jwtToken.SignedString(s.config.Secret)
}

func (s *JwtStorage) CreateTokenPair(ctx context.Context, claims AccessClaims) (TokenPair, error) {

	var (
		err error
	)

	tokens := TokenPair{}

	if tokens.Access, err = s.createToken(claims, time.Minute*15); err != nil {
		return tokens, err
	}
	if tokens.Refresh, err = s.createToken(claims, time.Hour*24*30); err != nil {
		return tokens, err
	}

	if err = s.recordToArchive(ctx, tokens); err != nil {
		return tokens, err
	}

	return tokens, nil
}
