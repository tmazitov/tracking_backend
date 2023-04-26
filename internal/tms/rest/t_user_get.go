package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type UserGetHandler struct {
	Jwt     jwt.JwtStorage
	Storage bl.Storage
	result  bl.R_GetUser
}

func (h *UserGetHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	userInfo, err := h.Storage.UserStorage().GetUserInfo(userPayload.UserId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result = bl.R_GetUser{
		ShortName: userInfo.ShortName.String,
		RoleID:    bl.UserRole(userInfo.RoleID.Int32),
	}

	core.SendResponse(200, h.result, ctx)
}
