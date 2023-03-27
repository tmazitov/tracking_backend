package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderGetHandler struct {
	Jwt     jwt.JwtStorage
	Storage bl.Storage
	result  bl.GetUser
}

func (h *OrderGetHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(403, "Unauthorized", err, ctx)
		return
	}

	userInfo, err := h.Storage.UserStorage().GetUserInfo(userPayload.UserId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result = bl.GetUser{
		TelNumber: userInfo.TelNumber.String,
		ShortName: userInfo.ShortName.String,
		RoleID:    userInfo.RoleID,
	}

	core.SendResponse(200, h.result, ctx)
}
