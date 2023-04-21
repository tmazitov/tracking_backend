package rest

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
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
		orderOwnerId         int64
		orderStatus          bl.OrderStatus
		updatedOrderPointsID []int64
		originalPointsID     []int64
		orderId              int64
		err                  error
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

	if len(h.input.Points) < 2 || len(h.input.Points) > 100 {
		core.ErrorLog(400, "Bad request", errors.New("count of points is not valid"), ctx)
		return
	}

	orderId, err = strconv.ParseInt(ctx.Param("orderId"), 10, 64)
	if err != nil || orderId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("upgrade order status: order_id is invalid"), ctx)
		return
	}

	// Check order owner
	orderOwnerId, err = h.Storage.OrderStorage().OrderGetOwnerID(orderId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}
	if int64(userPayload.UserId) != orderOwnerId {
		core.ErrorLog(403, "Forbidden", errors.New("invalid user id to update order"), ctx)
		return
	}

	// Check order status
	orderStatus, err = h.Storage.OrderStorage().OrderGetStatus(orderId)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}
	if orderStatus == bl.OrderCanceled || orderStatus == bl.OrderDone {
		core.ErrorLog(403, "Forbidden", errors.New("invalid user id to update order"), ctx)
		return
	}

	// Update basic info about the order (without points ID)
	if originalPointsID, err = h.Storage.OrderStorage().OrderGetPointsID(orderId); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	// Sort points on groups: to create, update, delete and make appropriate changes to points
	// Will return the array of points id which have relationship with current order
	if updatedOrderPointsID, err = h.updatePoints(orderId, originalPointsID, h.input.Points); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	// Update the points ID in the order record
	if err = h.Storage.OrderStorage().OrderUpdate(orderId, bl.DB_EditableOrder{
		StartAt:        h.input.StartAt,
		PointsID:       updatedOrderPointsID,
		Helpers:        h.input.Helpers,
		Comment:        h.input.Comment,
		IsFragileCargo: h.input.IsFragileCargo,
	}); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, nil, ctx)
}

func (h *OrderPutHandler) updatePoints(orderID int64, originalPointsID []int64, newPointsData []bl.Point) ([]int64, error) {
	var (
		err                  error
		idIsFound            bool
		pointsToUpdate       []bl.Point
		pointsToCreate       []bl.Point
		pointsToDelete       []int64
		createdPointsID      []int64
		updatedOrderPointsID []int64
	)

	// Find the new points which we need to create or update
	// !!! : all points with the extra ids is will be ignored
	for _, point := range newPointsData {
		idIsFound = false

		// If point don't have id -> to create array
		if point.ID == 0 {
			pointsToCreate = append(pointsToCreate, point)
			continue
		}

		for _, id := range originalPointsID {
			if point.ID == id {
				idIsFound = true
				break
			}
		}

		// If point have id and it is in the order points array -> update point data
		if idIsFound {
			pointsToUpdate = append(pointsToUpdate, point)
		}
	}

	// Find all points which we need to delete
	for _, originalPointId := range originalPointsID {
		idIsFound = false
		for _, pointToUpdate := range pointsToUpdate {
			if pointToUpdate.ID == originalPointId {
				idIsFound = true
				break
			}
		}

		if !idIsFound {
			pointsToDelete = append(pointsToDelete, originalPointId)
		}
	}

	// fmt.Println("points to create: ", pointsToCreate)
	// fmt.Println("points to update: ", pointsToUpdate)
	// fmt.Println("points to delete: ", pointsToDelete)

	if len(pointsToUpdate) > 0 {
		if updatedOrderPointsID, err = h.Storage.OrderStorage().PointsUpdate(pointsToUpdate); err != nil {
			fmt.Println("Update error")
			return nil, err
		}
	}
	if len(pointsToCreate) > 0 {
		if createdPointsID, err = h.Storage.OrderStorage().PointsCreate(orderID, pointsToCreate); err != nil {
			fmt.Println("Create error")
			return nil, err
		}
	}

	if len(pointsToDelete) > 0 {
		if err = h.Storage.OrderStorage().PointsDelete(pointsToDelete); err != nil {
			fmt.Println("Delete error")
			return nil, err
		}
	}

	updatedOrderPointsID = append(updatedOrderPointsID, createdPointsID...)
	return updatedOrderPointsID, nil
}
