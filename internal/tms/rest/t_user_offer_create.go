package rest

import (
	"errors"

	"github.com/gin-gonic/gin"
	bl "github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type UserOfferCreate struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	input   bl.UserJob
	result  struct {
		OfferId int `json:"offerId"`
	}
}

func (h *UserOfferCreate) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Base) {
		core.ErrorLog(403, "Forbidden", errors.New("offer create : invalid user role"), ctx)
		return
	}

	if err := ctx.ShouldBindJSON(&h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if h.input.JobType != uint8(bl.Manager) && h.input.JobType != uint8(bl.Worker) {
		core.ErrorLog(400, "Bad request", errors.New("offer create : invalid job type"), ctx)
		return
	}

	if h.input.JobExperience > 60 {
		core.ErrorLog(400, "Bad request", errors.New("offer create : invalid job experience"), ctx)
		return
	}

	var offerId int
	offerId, err = h.Storage.UserStorage().UserOfferGet(userPayload.UserId)
	if err != nil {
		core.ErrorLog(502, "Internal server error", err, ctx)
		return
	}
	if offerId != 0 {
		core.ErrorLog(403, "Forbidden", errors.New("offer create : offer already exists"), ctx)
		return
	}

	if h.result.OfferId, err = h.Storage.UserStorage().UserOfferCreate(userPayload.UserId, h.input); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(201, h.result, ctx)
}
