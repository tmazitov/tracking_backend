package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OrderDefaultVariables struct {
	Storage bl.Storage
	Jwt     jwt.JwtStorage
	result  struct {
		Prices []bl.DefaultPriceItems `json:"prices"`
	}
}

func (h *OrderDefaultVariables) Handle(ctx *gin.Context) {

	var (
		err error
	)

	_, err = h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if h.result.Prices, err = h.Storage.OrderStorage().OrderDefaultPrices(); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
