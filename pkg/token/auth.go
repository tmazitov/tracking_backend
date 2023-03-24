package token

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

func GetAuth(ctx *gin.Context) (string, error) {
	token := ctx.GetHeader("Authorization")

	// If no token
	if token == "" {
		return "", ErrInvalidToken
	}

	// If "Bearer" token
	if strings.Contains(token, "Bearer") {
		token = strings.ReplaceAll(token, "Bearer ", "")
	}

	return token, nil
}
