package rest

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/internal/tms/ws"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderTimeStartHandler struct {
	Storage bl.Storage
	Hub     *ws.Hub
	Jwt     jwt.JwtStorage
	query   struct {
		OrderId int64 `json:"orderId" bind:"required"`
	}
	result struct {
		StatusID    bl.OrderStatus `json:"statusId"`
		StartAtFact *time.Time     `json:"startAtFact"`
	}
}

func (h *OrderTimeStartHandler) Handle(ctx *gin.Context) {

	var (
		userPayload *jwt.AccessClaims
		err         error
	)

	userPayload, err = h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	h.query.OrderId, err = strconv.ParseInt(ctx.Param("orderId"), 10, 64)
	if err != nil || h.query.OrderId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("start order : order_id is invalid"), ctx)
		return
	}

	if userPayload.RoleId != int(bl.Worker) {
		core.ErrorLog(403, "Forbidden", errors.New("start order : user is not worker"), ctx)
		return
	}

	order, err := h.Storage.OrderStorage().OrderGet(h.query.OrderId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	if order.Worker.ID.Int64 != userPayload.UserId {
		core.ErrorLog(403, "Forbidden", errors.New("start order : user is not worker"), ctx)
		return
	}

	if h.result.StartAtFact, err = h.Storage.OrderStorage().OrderTimeStart(h.query.OrderId); err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}

	h.result.StatusID = 5
	go h.Hub.Broadcast(ws.OrderUpdateMessage{
		OrderId: h.query.OrderId,
		Type:    1,
		Data:    h.result,
	})

	h.Hub.UpdateStartAtFact(h.query.OrderId, h.result)
	core.SendResponse(200, h.result, ctx)
}
