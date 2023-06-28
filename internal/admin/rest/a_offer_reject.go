package rest

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OfferRejectHandler struct {
	Storage bl.Storage
	query   struct {
		OfferId int
	}
	result struct {
	}
}

func (h *OfferRejectHandler) Handle(ctx *gin.Context) {

	var err error

	h.query.OfferId, err = strconv.Atoi(ctx.Param("offerId"))
	if err != nil || h.query.OfferId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("offer accept : offer_id is invalid"), ctx)
		return
	}

	if err = h.Storage.UserStorage().OfferReject(h.query.OfferId); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
