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

type OrderSetWorkerHandler struct {
	Storage bl.Storage
	Hub     *ws.Hub
	Jwt     jwt.JwtStorage
	query   struct {
		OrderId int64 `json:"orderId" bind:"required"`
	}
	input struct {
		WorkerId int64 `json:"workerId" bind:"required"`
	}
	result struct {
		Worker   bl.R_GetUser   `json:"worker"`
		StatusId bl.OrderStatus `json:"statusId"`
	}
}

func (h *OrderSetWorkerHandler) Handle(ctx *gin.Context) {
	var (
		worker      *bl.DB_GetUser
		order       *bl.DB_Order
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
		core.ErrorLog(400, "Bad request", errors.New("set order worker : order_id is invalid"), ctx)
		return
	}

	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Manager) && userPayload.RoleId != int(bl.Admin) {
		core.ErrorLog(403, "Forbidden", errors.New("users can not set worker for the order"), ctx)
		return
	}

	if order, err = h.Storage.OrderStorage().OrderGet(h.query.OrderId); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if order.StatusID == int(bl.OrderCanceled) || order.StatusID == int(bl.OrderDone) {
		core.ErrorLog(403, "Forbidden", errors.New("status id is invalid for worker update of order"), ctx)
		return
	}

	if !order.Manager.ID.Valid || order.Manager.ID.Int64 != userPayload.UserId {
		core.ErrorLog(403, "Forbidden", errors.New("manager id is invalid for worker update of order"), ctx)
		return
	}

	if worker, err = h.Storage.OrderStorage().OrderUpdateWorker(h.query.OrderId, h.input.WorkerId); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result.Worker = bl.R_GetUser{
		ID:        worker.ID.Int64,
		ShortName: worker.ShortName.String,
		RoleID:    bl.UserRole(worker.RoleID.Int32),
	}

	h.result.StatusId = bl.OrderAccepted

	if order, err = h.Storage.OrderStorage().OrderGet(h.query.OrderId); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	var (
		updatedOrder bl.R_Order
		startAt      time.Time = order.StartAt.Time
		endAt        time.Time = order.EndAt.Time
	)

	updatedOrder = bl.R_Order{
		ID:                order.ID,
		Title:             order.Title,
		StartAt:           &startAt,
		EndAt:             &endAt,
		StartAtFact:       nil,
		EndAtFact:         nil,
		StatusID:          order.StatusID,
		Points:            order.Points,
		Comment:           order.Comment.String,
		IsRegularCustomer: order.IsRegularCustomer,
		Price: &bl.R_OrderBill{
			CarTypeID:      order.Bill.CarTypeID,
			HelperCount:    uint(order.Bill.HelperCount.Int16),
			HelperHours:    uint(order.Bill.HelperHours.Int16),
			IsFragileCargo: order.Bill.IsFragileCargo,
		},
	}

	owner := bl.R_GetUser{
		ID:        order.Owner.ID.Int64,
		ShortName: order.Owner.ShortName.String,
		RoleID:    bl.UserRole(order.Owner.RoleID.Int32),
	}
	updatedOrder.Owner = &owner

	if order.Worker.ID.Valid {
		var worker bl.R_GetUser = bl.R_GetUser{
			ID:        order.Worker.ID.Int64,
			ShortName: order.Worker.ShortName.String,
			RoleID:    bl.UserRole(order.Worker.RoleID.Int32),
		}
		updatedOrder.Worker = &worker
	}

	if order.Manager.ID.Valid {
		var manager bl.R_GetUser = bl.R_GetUser{
			ID:        order.Manager.ID.Int64,
			ShortName: order.Manager.ShortName.String,
			RoleID:    bl.UserRole(order.Manager.RoleID.Int32),
		}
		updatedOrder.Manager = &manager
	}

	if err = h.Hub.UpdateWorker(ctx, &updatedOrder, h.result); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
