package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CheckClaims struct {
	Ip        string    `json:"ip"`
	Email     string    `json:"email"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (c *CheckClaims) Valid() error {
	if time.Now().After(c.ExpiredAt) {
		return ErrInvalidToken
	}
	return nil
}

func (j *JwtStorage) NewCheckToken(claims *CheckClaims) (string, error) {

	now := time.Now()

	claims.IssuedAt = now
	claims.ExpiredAt = now.Add(5 * time.Minute)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(j.config.Secret)
}
