package rest

import (
	"errors"
	"time"

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

	var (
		startAt time.Time = order.StartAt.Time
		endAt   time.Time = order.EndAt.Time
	)

	h.result = bl.R_Order{
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
	h.result.Owner = &owner

	if order.Manager.ID.Valid {
		var manager bl.R_GetUser = bl.R_GetUser{
			ID:        order.Manager.ID.Int64,
			ShortName: order.Manager.ShortName.String,
			RoleID:    bl.UserRole(order.Manager.RoleID.Int32),
		}
		h.result.Manager = &manager
	}

	core.SendResponse(201, h.result, ctx)
}
