package rest

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type UserHolidayListByDate struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	query   struct {
		Date time.Time
	}
	result []int64
}

func (h *UserHolidayListByDate) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Manager) && userPayload.RoleId != int(bl.Admin) {
		core.ErrorLog(403, "Forbidden", errors.New("worker can not to create the order"), ctx)
		return
	}

	params := ctx.Request.URL.Query()
	h.query.Date, err = time.Parse("2006-01-02", params.Get("d"))
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if h.result, err = h.Storage.UserStorage().UserHolidayListByDate(h.query.Date); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(200, h.result, ctx)
}
