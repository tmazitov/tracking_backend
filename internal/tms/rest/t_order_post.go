package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type AddOrderHandler struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	input   struct {
		StartAt        time.Time  `json:"startAt" validate:"max=32"`
		Points         []bl.Point `json:"points"`
		Helpers        uint8      `json:"helpers,omitempty"`
		Comment        string     `json:"comment,omitempty" validate:"max=256"`
		IsFragileCargo bool       `json:"isFragileCargo,omitempty"`
	}
	result struct {
	}
}

func (h *AddOrderHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(403, "Unauthorized", err, ctx)
		return
	}

	if err := ctx.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	pointIDs, err := h.Storage.OrderStorage().InsertPoint(h.input.Points)
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
	if err = h.Storage.OrderStorage().InsertOrder(newOrder); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(201, h.result, ctx)
}
