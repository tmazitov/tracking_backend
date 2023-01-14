package rest

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/app/bl"
)

type AddOrderHandler struct {
	Storage bl.Storage
	input   struct {
		StartAt time.Time  `json:"startAt" validate:"max=32"`
		EndAt   time.Time  `json:"endAt" validate:"max=32"`
		Points  []bl.Point `json:"points"`
		Helpers uint8      `json:"helpers,omitempty"`
		Comment string     `json:"comment,omitempty" validate:"max=256"`
	}
	result struct {
	}
}

func (h *AddOrderHandler) Input() interface{} {
	return &h.input
}
func (h *AddOrderHandler) Result() interface{} {
	return &h.result
}

func (h *AddOrderHandler) Handle(c *gin.Context) {

	if err := c.BindJSON(&h.input); err != nil {
		ErrorLog(400, "Bad request", err, c)
		return
	}

	pointIDs, err := h.Storage.OrderStorage().InsertPoint(h.input.Points)
	if err != nil {
		ErrorLog(500, "Internal server error", err, c)
		return
	}

	newOrder := bl.CreateOrder{
		StartAt:  h.input.StartAt,
		EndAt:    h.input.EndAt,
		PointsID: pointIDs,
		Helpers:  h.input.Helpers,
		Comment:  h.input.Comment,
	}
	if err = h.Storage.OrderStorage().InsertOrder(newOrder); err != nil {
		ErrorLog(500, "Internal server error", err, c)
		return
	}

	SendResponse(201, h.result, c)
}
