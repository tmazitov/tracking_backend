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

	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	pointIDs, err := h.Storage.OrderStorage().CreatePoints(h.input.Points)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	newOrder := bl.CreateOrder{
		StartAt:        h.input.StartAt,
		OwnerID:        userPayload.UserId,
		PointsID:       pointIDs,
		Helpers:        h.input.Helpers,
		Comment:        h.input.Comment,
		IsFragileCargo: h.input.IsFragileCargo,
	}
	isManager := userPayload.RoleId == int(bl.Manager)
	if err = h.Storage.OrderStorage().CreateOrder(newOrder, isManager); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(201, nil, ctx)
}
