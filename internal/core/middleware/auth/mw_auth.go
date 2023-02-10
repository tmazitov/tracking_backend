package auth

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/core/jwt"
	core "github.com/tmazitov/tracking_backend.git/internal/core/request"
)

var (
	ErrUnauthorized = errors.New("user is not authorized")
)

type Middleware struct {
	Jwt *jwt.JwtStorage
}

func (mw *Middleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")

		// If no token
		if token == "" {
			core.ErrorLog(401, "unauthorized ", ErrUnauthorized, ctx)
			return
		}

		// If "Bearer" token
		if strings.Contains(token, "Bearer") {
			token = strings.ReplaceAll(token, "Bearer ", "")
		}

		// Validate accesss token
		if _, err := mw.Jwt.ValidateAccess(ctx, token); err != nil {
			core.ErrorLog(401, "unauthorized", err, ctx)
			return
		}

		// after request

		ctx.Next()
	}
}
