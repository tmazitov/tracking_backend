package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/admin/bl"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type OfferListHandler struct {
	Storage bl.Storage
	result  struct {
		Offers []bl.R_UserOffer `json:"offers"`
	}
}

func (h *OfferListHandler) Handle(ctx *gin.Context) {
	var (
		offersRaw []bl.DB_UserOffer
		err       error
	)
	if offersRaw, err = h.Storage.UserStorage().UserOfferList(); err != nil {
		core.ErrorLog(502, "Internal server error", err, ctx)
		return
	}

	h.result.Offers = []bl.R_UserOffer{}
	for _, offerRaw := range offersRaw {
		h.result.Offers = append(h.result.Offers, *offerRaw.ToReal())
	}

	core.SendResponse(200, h.result, ctx)
}
