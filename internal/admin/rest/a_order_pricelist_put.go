package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderPriceListPutHandler struct {
	Storage bl.Storage
	input   bl.OrderPriceList
}

func (h *OrderPriceListPutHandler) Handle(ctx *gin.Context) {
	if err := ctx.ShouldBindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if err := h.Storage.OrderStorage().OrderPricelistUpdate(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(200, nil, ctx)
}
