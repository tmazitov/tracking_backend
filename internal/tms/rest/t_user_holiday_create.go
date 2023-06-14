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

type UserHolidayCreateHandler struct {
	Storage bl.Storage
	Jwt     *jwt.JwtStorage
	query   struct {
		WorkerId int64
		Date     time.Time `json:"date"`
	}
}

func (h *UserHolidayCreateHandler) Handle(ctx *gin.Context) {
	userPayload, err := h.Jwt.ValidateAccess(ctx)
	if err != nil {
		core.ErrorLog(401, "Unauthorized", err, ctx)
		return
	}

	if userPayload.RoleId != int(bl.Admin) {
		core.ErrorLog(403, "Forbidden", errors.New("worker can not to create the order"), ctx)
		return
	}

	h.query.WorkerId, err = strconv.ParseInt(ctx.Param("userId"), 10, 64)
	if err != nil || h.query.WorkerId <= 0 {
		core.ErrorLog(400, "Bad request", errors.New("set worker holiday : user_id is invalid"), ctx)
		return
	}

	params := ctx.Request.URL.Query()
	h.query.Date, err = time.Parse("2006-01-02", params.Get("d"))
	if err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	if err = h.Storage.UserStorage().UserHolidayCreate(h.query.WorkerId, userPayload.UserId, &h.query.Date); err != nil {
		core.ErrorLog(400, "Bad request", err, ctx)
		return
	}

	core.SendResponse(201, nil, ctx)
}
