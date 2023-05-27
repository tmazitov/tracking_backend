package rest

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	bl "github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/internal/tms/ws"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderTimeEndHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	Hub     *ws.Hub
	query   struct {
		OrderId int64 `json:"orderId" bind:"required"`
	}
	result struct {
		StatusID  bl.OrderStatus `json:"statusId"`
		EndAtFact *time.Time     `json:"endAtFact"`
	}
}

func (h *OrderTimeEndHandler) Handle(ctx *gin.Context) {

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

	if h.result.EndAtFact, err = h.Storage.OrderStorage().OrderTimeEnd(h.query.OrderId); err != nil {
		core.ErrorLog(500, "Internal Server error", err, ctx)
		return
	}
	h.result.StatusID = 1

	h.Hub.UpdateEndAtFact(h.query.OrderId, h.result)
	core.SendResponse(200, h.result, ctx)
}
