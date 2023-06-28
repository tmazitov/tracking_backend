package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type UserOfferGet struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	result  struct {
		OfferId int `json:"offerId,omitempty"`
	}
}

func (h *UserOfferGet) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Base) {
		core.ErrorLog(403, "Forbidden", errors.New("offer get : invalid user role"), ctx)
		return
	}

	if h.result.OfferId, err = h.Storage.UserStorage().UserOfferGet(userPayload.UserId); err != nil {
		core.ErrorLog(502, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
