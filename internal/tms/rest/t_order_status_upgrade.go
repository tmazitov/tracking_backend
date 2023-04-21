package rest

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderStatusUpgradeHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	result  struct {
		StatusID int `json:"statusId"`
	}
}

func (h *OrderStatusUpgradeHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Admin) && userPayload.RoleId != int(bl.Manager) {
		core.ErrorLog(403, "Forbidden", errors.New("worker or user can not update status of order"), ctx)
		return
	}

	orderId, err := strconv.ParseInt(ctx.Param("orderId"), 10, 64)
	if err != nil || orderId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("upgrade order status: order_id is invalid"), ctx)
		return
	}

	managerId, err := h.Storage.OrderStorage().OrderGetManagerID(orderId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	if managerId != int64(userPayload.UserId) {
		core.ErrorLog(403, "Forbidden", errors.New("invalid manager id"), ctx)
		return
	}

	if h.result.StatusID, err = h.Storage.OrderStorage().OrderStatusUpgrade(orderId); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
