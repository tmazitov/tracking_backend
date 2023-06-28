package rest

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OfferAcceptHandler struct {
	Storage bl.Storage
	query   struct {
		OfferId int
	}
	result bl.R_User
}

func (h *OfferAcceptHandler) Handle(ctx *gin.Context) {

	var (
		err     error
		userRaw *bl.DB_User
	)

	h.query.OfferId, err = strconv.Atoi(ctx.Param("offerId"))
	if err != nil || h.query.OfferId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("offer accept : offer_id is invalid"), ctx)
		return
	}

	if userRaw, err = h.Storage.UserStorage().OfferAccept(h.query.OfferId); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	h.result = *userRaw.ToReal()
	core.SendResponse(200, h.result, ctx)
}
