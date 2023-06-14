package rest

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tmazitov/tracking_backend.git/internal/tms/bl"
	"github.com/tmazitov/tracking_backend.git/pkg/jwt"
	core "github.com/tmazitov/tracking_backend.git/pkg/request"
)

type UserHolidayDeleteHandler struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	query   struct {
		Date     time.Time
		WorkerId int64
	}
}

func (h *UserHolidayDeleteHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Admin) {
		core.ErrorLog(403, "Forbidden", errors.New("only manager and admin can delete holiday"), ctx)
		return
	}

	h.query.WorkerId, err = strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil || h.query.WorkerId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("delete worker holiday : user_id is invalid"), ctx)
		return
	}

	params := ctx.Request.URL.Query()
	h.query.Date, err = time.Parse("2006-01-02", params.Get("d"))
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	var holiday *bl.UserHoliday
	if holiday, err = h.Storage.UserStorage().UserHolidayGet(h.query.WorkerId, &h.query.Date); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if holiday.AuthorId == 0 {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if err = h.Storage.UserStorage().UserHolidayDelete(h.query.WorkerId, &h.query.Date); err != nil {
		core.ErrorLog(500, "Internal server error", err, ctx)
		return
	}

	core.SendResponse(200, nil, ctx)
}
