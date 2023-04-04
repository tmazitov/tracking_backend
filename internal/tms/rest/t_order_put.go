package rest

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	bl "github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderPutHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	input   bl.R_EditableOrder
}

func (h *OrderPutHandler) Handle(ctx *gin.Context) {

	var (
		pointsID pq.Int64Array
		orderId  int
		err      error
	)

	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId == int(bl.Worker) {
		core.ErrorLog(403, "Forbidden", errors.New("worker can not to create the order"), ctx)
		return
	}

	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	orderId, err = strconv.Atoi(ctx.Param("orderId"))
	if err != nil || orderId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("upgrade order status: order_id is invalid"), ctx)
		return
	}

	if pointsID, err = h.Storage.OrderStorage().OrderUpdateMainInfo(orderId, h.input); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	if len(pointsID) == len(h.input.Points) {
		// update one by one
	} else if len(pointsID) < len(h.input.Points) {
		// update points and insert new
	} else if len(pointsID) > len(h.input.Points) {
		// update points and update old
	}

	core.SendResponse(200, nil, ctx)
}
