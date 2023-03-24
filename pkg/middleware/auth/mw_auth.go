package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type Middleware struct {
	Jwt *jwt.JwtStorage
}

func (mw *Middleware) Handle() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Validate accesss token
		if _, err := mw.Jwt.ValidateAccess(ctx); err != nil {
			core.ErrorLog(401, "unauthorized", err, ctx)
			return
		}

		// after request
		ctx.Next()
	}
}
