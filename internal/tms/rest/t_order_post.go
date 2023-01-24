package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	core "github.com/tmazitov/tracking_backend.git/internal/core/request"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
)

type AddOrderHandler struct {
	Storage bl.Storage
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

func (h *AddOrderHandler) Handle(c *gin.Context) {

	if err := c.BindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, c)
		return
	}

	pointIDs, err := h.Storage.OrderStorage().InsertPoint(h.input.Points)
	if err != nil {
		core.ErrorLog(500, "Internal server error", err, c)
		return
	}

	newOrder := bl.CreateOrder{
		StartAt:        h.input.StartAt,
		PointsID:       pointIDs,
		Helpers:        h.input.Helpers,
		Comment:        h.input.Comment,
		IsFragileCargo: h.input.IsFragileCargo,
	}
	if err = h.Storage.OrderStorage().InsertOrder(newOrder); err != nil {
		core.ErrorLog(500, "Internal server error", err, c)
		return
	}

	core.SendResponse(201, h.result, c)
}
