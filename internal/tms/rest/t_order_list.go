package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderListHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	// input   struct {
	// }
	result []bl.R_OrderListItem
}

func (h *OrderListHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if h.result, err = getOrderList(userPayload.UserId, userPayload.RoleId, h.Storage); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}

func getOrderList(userId int, roleId int, storage bl.Storage) ([]bl.R_OrderListItem, error) {
	var (
		orders []bl.DB_OrderListItem
		result []bl.R_OrderListItem
		err    error
	)

	if orders, err = storage.OrderStorage().OrderList(userId, roleId); err != nil {
		return result, err
	}

	for _, order := range orders {
		points, err := storage.OrderStorage().PointsGet(order.PointsID)
		if err != nil {
			return result, err
		}

		result_item := bl.R_OrderListItem{
			ID:                order.ID,
			Title:             order.Title,
			StartAt:           order.StartAt,
			EndAt:             order.EndAt.Time,
			StatusID:          order.StatusID,
			Points:            points,
			OwnerID:           order.OwnerID,
			WorkerID:          int(order.WorkerID.Int64),
			ManagerID:         int(order.ManagerID.Int64),
			Helpers:           uint8(order.Helpers.Int16),
			Comment:           order.Comment.String,
			IsFragileCargo:    order.IsFragileCargo,
			IsRegularCustomer: order.IsRegularCustomer,
		}
		result = append(result, result_item)
	}

	return result, err
}
