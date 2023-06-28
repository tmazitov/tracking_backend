package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type AdminRoleMiddleware struct {
	Jwt *jwt.JwtStorage
}

func (h *AdminRoleMiddleware) Handle(ctx *gin.Context) {
	var (
		userPayload *jwt.AccessClaims
		err         error
	)

	// Validate access token
	if userPayload, err = h.Jwt.ValidateAccess(ctx); err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		ctx.Abort()
		return
	}

	// Validate access level
	if userPayload.RoleId != int(bl.Admin) {
		core.ErrorLog(403, "Forbidden", err, ctx)
		ctx.Abort()
		return
	}

	ctx.Next()
}
