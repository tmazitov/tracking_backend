package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderCreateHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	input   bl.R_CreatableOrder
	result  bl.R_Order
}

func (h *OrderCreateHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId == int(bl.Worker) {
		core.ErrorLog(403, "Forbidden", errors.New("worker can not to create the order"), ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if len(h.input.Points) < 2 || len(h.input.Points) > 100 {
		core.ErrorLog(400, "Bad request", errors.New("count of points is not valid"), ctx)
		return
	}

	if h.input.OrderType > 7 {
		core.ErrorLog(400, "Bad request", errors.New("order type is greater than 6"), ctx)
		return
	}

	var (
		newOrder bl.CreateOrder
	)

	if h.input.Title == "" {
		h.input.Title = h.input.Points[0].Title
	}

	newOrder = bl.CreateOrder{
		OwnerID:           userPayload.UserId,
		WorkerID:          h.input.WorkerID,
		StartAt:           h.input.StartAt,
		EndAt:             h.input.EndAt,
		Title:             h.input.Title,
		Points:            h.input.Points,
		OrderType:         h.input.OrderType,
		Comment:           h.input.Comment,
		IsRegularCustomer: h.input.IsRegularCustomer,
	}

	var orderID int64
	if orderID, err = h.Storage.OrderStorage().CreateOrder(newOrder, bl.UserRole(userPayload.RoleId)); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	if err = h.Storage.OrderStorage().OrderBillUpdatePrice(orderID, h.input.Price); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	var order *bl.DB_Order
	if order, err = h.Storage.OrderStorage().OrderGet(orderID); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	h.result = *order.ToReal()

	core.SendResponse(201, h.result, ctx)
}
