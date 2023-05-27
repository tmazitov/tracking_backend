package jwt

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
)

func (s *JwtStorage) AccessIsExists(ctx context.Context, token string) error {
	return s.isExists(ctx, "acc:", token)
}

func (s *JwtStorage) ValidateAccess(ctx *gin.Context) (*AccessClaims, error) {

	token, err := GetAccessFromParams(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	return s.ValidateAccessToken(ctx, token)
}

func (s *JwtStorage) ValidateAccessToken(ctx context.Context, token string) (*AccessClaims, error) {

	// Check if token is exists
	if err := s.AccessIsExists(ctx, token); err != nil {
		return nil, err
	}

	return s.verifyToken(ctx, token)
}

func GetAccessFromParams(ctx *gin.Context) (string, error) {
	var (
		token string
		err   error
	)

	token = ctx.GetHeader("Authorization")

	// If no token
	if token == "" {
		return "", ErrUnauthorized
	}

	// If "Bearer" token
	if strings.Contains(token, "Bearer") {
		token = strings.ReplaceAll(token, "Bearer ", "")
	}

	return token, err
}
